/*

This program counts words in text files, given as command line parameters.

There is no concurrency, nor any channels

Author: Jamie Skipworth

*/
package main
import (
	"fmt"
	"os"
	"time"
	"./demo1"
	"log"
	)

var logger *log.Logger

// Count some words!
func main( ) {

	// Set the format of the logger
	log.SetFlags( log.Ldate|log.Lmicroseconds )

	// Make sure we have at least one cmd line parameter
	if( len(os.Args) < 2 ){
		log.Output( 2, fmt.Sprintf( "Expected at least one filename." ) )
		os.Exit(1)
	}

	// Take a slice of the command line. Arg[0] is the program name.
	filenames := os.Args[1:]

	// Set the start time so we can measure performance.
	timeStart := time.Now()

	// Iterate through files from the command line arguments
	for i := 0; i < len( filenames ); i++ {

		// Call the wordCounter. Discard bytesRead (1st return param), but keep the word count (2nd)
		_, w := demo1.WordCounter( filenames[i] )
		log.Output( 2, fmt.Sprintf( "%d\twords in %s", w, filenames[i] ) )
	}

	// Show total elapsed time.
	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "Total elapsed time %s", timeEnd.Sub( timeStart )) )


}
