package transport

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	RetryCount = 5
	Backoff    = 100
)

var (
	retriableErrors = []codes.Code{codes.Unavailable, codes.DataLoss, codes.Unknown, codes.NotFound, codes.Aborted}
)
var (
	kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than onIf a client pings more than oce every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
)

func NewGrpcConn(addr string, receiveSize int) *grpc.ClientConn {
	var (
		err  error
		m    = 1024 * 1024
		conn *grpc.ClientConn
	)
	opts := []grpc_retry.CallOption{
		grpc_retry.WithCodes(retriableErrors...),
		grpc_retry.WithMax(RetryCount),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(Backoff * time.Millisecond)),
	}
	conn, err = grpc.Dial(addr,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(receiveSize*m)),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_prometheus.UnaryClientInterceptor,
				grpc_retry.UnaryClientInterceptor(opts...))),
	)
	if err != nil {
		log.Panic("grpc conn error ==>", err)
	}
	if state := conn.GetState(); state == connectivity.Shutdown {
		log.Panic("grpc conn shutdown ==>")
	}
	return conn
}

func NewGrpcServer() (rs *grpc.Server) {

	var zapLogger, _ = zap.NewProduction()

	rcOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}

	zpOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(func(c codes.Code) zapcore.Level {
			level := ClientCodeToLevel(c)
			return level
		}),
		grpc_zap.WithDurationField(grpc_zap.DurationToDurationField),
	}
	grpc_zap.ReplaceGrpcLoggerV2(zapLogger)
	return grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(rcOpts...),
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_zap.UnaryServerInterceptor(zapLogger, zpOpts...),
			)),
	)

}

func ClientCodeToLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zap.DebugLevel
	case codes.Canceled:
		return zap.DebugLevel
	case codes.Unknown:
		return zap.InfoLevel
	case codes.InvalidArgument:
		return zap.DebugLevel
	case codes.DeadlineExceeded:
		return zap.InfoLevel
	case codes.NotFound:
		return zap.DebugLevel
	case codes.AlreadyExists:
		return zap.DebugLevel
	case codes.PermissionDenied:
		return zap.InfoLevel
	case codes.Unauthenticated:
		return zap.InfoLevel // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return zap.DebugLevel
	case codes.FailedPrecondition:
		return zap.DebugLevel
	case codes.Aborted:
		return zap.DebugLevel
	case codes.OutOfRange:
		return zap.DebugLevel
	case codes.Unimplemented:
		return zap.WarnLevel
	case codes.Internal:
		return zap.WarnLevel
	case codes.Unavailable:
		return zap.WarnLevel
	case codes.DataLoss:
		return zap.WarnLevel
	default:
		return zap.InfoLevel
	}
}
