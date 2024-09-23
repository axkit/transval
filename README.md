# transval - Transition Validator

[![Build Status](https://github.com/axkit/transval/actions/workflows/go.yml/badge.svg)](https://github.com/axkit/transval/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axkit/transval)](https://goreportcard.com/report/github.com/axkit/transval)
[![GoDoc](https://pkg.go.dev/badge/github.com/axkit/transval)](https://pkg.go.dev/github.com/axkit/transval)
[![Coverage Status](https://coveralls.io/repos/github/axkit/transval/badge.svg?branch=main)](https://coveralls.io/github/axkit/transval?branch=main)

`transval` is a Go package that provides functionality to manage and validate transitions between different states. This can be useful in scenarios like workflow management, state machines, or any system that involves transitioning from one state to another.

## Features

- Define transition rules between states.
- Check if a transition between two states is valid.
- Retrieve all possible transitions from a given state.
- Easily add, update, or delete named transition rules.

## Installation

To install the package, use the following command:

```bash
go get github.com/axkit/transval
```

# Usage

## Creating a New TransVal Instance

To create a new instance of the TransVal structure, use the New() function:
```go
import "transval"

trans := transval.New()
```

## Adding Transition Rules

You can define state transitions by using the Set method. The transition rules should follow the format:
```
<from_state> => <to_state1>,<to_state2>,...;<from_state2> => <to_state3>,...
```

For example:
```go
err := trans.Set("myRules", "1 => 2,3; 2 => 4")
if err != nil {
    // handle error
}
```

In this example, state 1 can transition to states 2 and 3, while state 2 can transition to state 4.

## Validating a Transition

To check if a transition between two states is valid, use the IsTransitionValid method:

```go
isValid := trans.IsTransitionValid("myRules", 1, 2)  // returns true
isValid = trans.IsTransitionValid("myRules", 1, 4)    // returns false
```

## Retrieving Allowed Transitions

To get a list of states that can be transitioned to from a given state, use the AllowedTo method:

```go
allowedStates := trans.AllowedTo("myRules", 1)  // returns []int{2, 3}
```

## Retrieving All Transitions

To get the entire set of transitions for a given rule set, use the Transitions method:
```go
allTransitions := trans.Transitions("myRules")
/*
    allTransitions = map[int][]int{
        1: {2, 3},
        2: {4},
    }
*/
```

## Deleting Transition Rules

To remove a set of transition rules, use the Del method:
```go
trans.Del("myRules")
```

## Errors

The package provides the following error values for common error conditions:

	•	ErrWrongInput: Indicates that the input format for the transition rules is incorrect.
	•	ErrEmptyInput: Indicates that the input string is empty.
	•	ErrTargetEmpty: Indicates that a target state value is missing in the transition rule.

Example:
```go
package main

import (
    "fmt"
    "log"
    "transval"
)

func main() {
    // Create a new TransVal instance
    trans := transval.New()

    // Define some transition rules
    err := trans.Set("workflow", "1 => 2,3; 2 => 4")
    if err != nil {
        log.Fatal(err)
    }

    // Check if a transition is valid
    fmt.Println(trans.IsTransitionValid("workflow", 1, 2))  // true
    fmt.Println(trans.IsTransitionValid("workflow", 1, 4))  // false

    // Get allowed transitions from state 1
    fmt.Println(trans.AllowedTo("workflow", 1))  // [2 3]

    // Get all transitions
    fmt.Println(trans.Transitions("workflow"))
}
```

# License

This package is available under the MIT License. See the LICENSE file for more information.

# Contributions

Contributions are welcome! Feel free to open issues or submit pull requests to improve the package.