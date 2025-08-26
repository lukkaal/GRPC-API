package shutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
)

// only close the httpserver(could add more rpc service)
func GracefullyShutdown(server *http.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT,
		syscall.SIGTERM)
	// release chan os.Signal
	defer signal.Stop(sigCh)

	// wait until the channel to notify
	done := <-sigCh

	// logrus.WithField("user", "luka").Info("login success")
	// time="2025-08-26 11:12:33" level=info msg="login success" user=luka
	logger.GinloggerObj.WithField("signal", done).
		Info("shutdown signal detected")

	// set context in case of jam
	ctx, cancel := context.WithTimeout(
		context.Background(), 15*time.Second)

	// incur ctx.Done() to release resource
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.GinloggerObj.Fatalln("closing http server gracefully failed:", err)
	}
}
