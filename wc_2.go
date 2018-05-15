/*

This program counts words in text files, given as command line parameters.

There is concurrency using goroutines, but no channels

Author: Jamie Skipworth

*/
package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"./demo2"
	"log"
	)


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

	// **** WaitGroups ****
	// A WaitGroup is useful when you have a single channel shared amongst
	// goroutines. How do you know when all the goroutines are done with it?
	//
	// When you start a goroutine that uses a shared channel, you Add() to the WaitGroup.
	// Once processing is finished, it calls Done(). Somewhere you will have to Wait() until
	// all the Add()s have a corresponding Done(). Then the Wait() will unblock and you can
	// close your shared channel.
	wg := sync.WaitGroup{}

	// Iterate through files from the command line arguments
	for i := 0; i < len( filenames ); i++ {	

		// **** WaitGroups ***
		// How do we know when the gorountines are finished counting words? We must use a WaitGroup.
		// We Add() goroutines we've started, which then call Done() when complete.
		wg.Add( 1 )

		// Execute wordCounter. We pass it he filename we're processing and a WaitGroup to send Done() to
		go demo2.WordCounter( filenames[i], &wg )
	}

	// This call blocks until all the goroutines are Done()
	wg.Wait()

	// Show total elapsed time.
	timeEnd := time.Now()
	log.Output( 2, fmt.Sprintf( "Total elapsed time %s", timeEnd.Sub( timeStart )) )
		
	
}

