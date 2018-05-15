package util
import(
	"fmt"
	"os"
	"time"
	)

// Writes text to a file descriptor. Nice for consistent logging.
func LogLine( fd *os.File, message string ) {
	fmt.Fprintf( fd, "%s - %s\n", time.Now().Format("20060106 15:04:05.000"), message )
}
