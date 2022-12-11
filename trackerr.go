// Package trackerr aims to facilitate creation of referenceable errors and
// elegant stack traces.
//
//		I crafted this package in reponse to my frustration in trying to
//		decypher Go's printed error stack traces and the challenge of reliably
//		asserting specific error types while testing.
//
//		Many programmers assert using error messages but I've found this to be
//		unreliable and leave me less than confident. trackerr attempts to rectify
//		this by assigning tracked errors there own unique identifiers which can
//		be checked using errors.Is or one of trackerr's utility functions.
//
//		Paulio
//
// The recommended way to create errors is via the Track, Checkpoint,
// Untracked, and Wrap package functions. It is not recommended to create
// trackable errors after initialisation but Realms exist for such cases.
//
// It is also recommended to call the Initialised function from an init
// function in package main to prevent creation of trackable errors after
// program initialisation.
package trackerr
