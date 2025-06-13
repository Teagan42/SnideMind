// File: utils/timer_test.go
package utils

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

// captureOutput is a helper function to capture log output.
func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil) // Reset output to default (stderr)
	}()
	f()
	return buf.String()
}

// TestTimeFunc tests the TimeFunc wrapper.
func TestTimeFunc(t *testing.T) {
	testName := "MySimpleFunc"
	expectedResult := 42

	// Define a simple function to be timed
	simpleFunc := func(input int) int {
		time.Sleep(10 * time.Millisecond) // Simulate some work
		return input * 2
	}

	// Wrap the function
	timedFunc := TimeFunc(testName, simpleFunc)

	// Capture output and run the timed function
	output := captureOutput(func() {
		result := timedFunc(21)
		// Verify the result
		if result != expectedResult {
			t.Errorf("TimeFunc(%s) returned %v, expected %v", testName, result, expectedResult)
		}
	})

	// Verify the log output
	if !strings.Contains(output, "[Timer]") {
		t.Errorf("Log output missing [Timer] prefix: %s", output)
	}
	if !strings.Contains(output, testName) {
		t.Errorf("Log output missing function name '%s': %s", testName, output)
	}
	if !strings.Contains(output, "took") {
		t.Errorf("Log output missing 'took' keyword: %s", output)
	}
	// Check if a duration is present (e.g., "10ms") - this is a basic check
	// A more robust check might parse the duration, but this is sufficient for a unit test
	if !strings.Contains(output, "ms") && !strings.Contains(output, "µs") && !strings.Contains(output, "ns") {
		t.Errorf("Log output missing duration indicator (ms, µs, ns): %s", output)
	}
}

// TestTimeFuncWithErr tests the TimeFuncWithErr wrapper.
func TestTimeFuncWithErr(t *testing.T) {
	testName := "MyFuncWithErr"
	expectedResult := "success"
	expectedErr := "something went wrong"

	// Test case 1: Function returns value and no error
	t.Run("NoError", func(t *testing.T) {
		funcNoError := func(input string) (string, error) {
			time.Sleep(10 * time.Millisecond)
			return input, nil
		}

		timedFunc := TimeFuncWithErr(testName, funcNoError)

		output := captureOutput(func() {
			result, err := timedFunc(expectedResult)
			if result != expectedResult {
				t.Errorf("TimeFuncWithErr(%s) returned result %v, expected %v", testName, result, expectedResult)
			}
			if err != nil {
				t.Errorf("TimeFuncWithErr(%s) returned unexpected error: %v", testName, err)
			}
		})

		if !strings.Contains(output, testName) || !strings.Contains(output, "took") {
			t.Errorf("Log output verification failed for NoError case: %s", output)
		}
	})

	// Test case 2: Function returns value and an error
	t.Run("WithError", func(t *testing.T) {
		funcWithError := func(input string) (string, error) {
			time.Sleep(10 * time.Millisecond)
			return input, fmt.Errorf(expectedErr)
		}

		timedFunc := TimeFuncWithErr(testName, funcWithError)

		output := captureOutput(func() {
			result, err := timedFunc(expectedResult)
			if result != expectedResult {
				t.Errorf("TimeFuncWithErr(%s) returned result %v, expected %v", testName, result, expectedResult)
			}
			if err == nil || err.Error() != expectedErr {
				t.Errorf("TimeFuncWithErr(%s) returned error %v, expected %v", testName, err, expectedErr)
			}
		})

		if !strings.Contains(output, testName) || !strings.Contains(output, "took") {
			t.Errorf("Log output verification failed for WithError case: %s", output)
		}
	})
}

// TestTimeFunc2WithErr tests the TimeFunc2WithErr wrapper.
func TestTimeFunc2WithErr(t *testing.T) {
	testName := "MyFunc2WithErr"
	expectedResult := 100
	expectedErr := "division by zero"

	// Test case 1: Function returns value and no error
	t.Run("NoError", func(t *testing.T) {
		funcNoError := func(a, b int) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return a * b, nil
		}

		timedFunc := TimeFunc2WithErr(testName, funcNoError)

		output := captureOutput(func() {
			result, err := timedFunc(10, 10)
			if result != expectedResult {
				t.Errorf("TimeFunc2WithErr(%s) returned result %v, expected %v", testName, result, expectedResult)
			}
			if err != nil {
				t.Errorf("TimeFunc2WithErr(%s) returned unexpected error: %v", testName, err)
			}
		})

		if !strings.Contains(output, "[Timer] Timing "+testName) {
			t.Errorf("Log output missing initial timing message: %s", output)
		}
		if !strings.Contains(output, "[Timer] "+testName+" took") {
			t.Errorf("Log output missing final took message: %s", output)
		}
	})

	// Test case 2: Function returns value and an error
	t.Run("WithError", func(t *testing.T) {
		funcWithError := func(a, b int) (int, error) {
			time.Sleep(10 * time.Millisecond)
			if b == 0 {
				return 0, fmt.Errorf(expectedErr)
			}
			return a / b, nil
		}

		timedFunc := TimeFunc2WithErr(testName, funcWithError)

		output := captureOutput(func() {
			_, err := timedFunc(10, 0)
			// We don't care about the result value when there's an error in this specific function
			// if result != 0 {
			// 	t.Errorf("TimeFunc2WithErr(%s) returned result %v, expected %v on error", testName, result, 0)
			// }
			if err == nil || err.Error() != expectedErr {
				t.Errorf("TimeFunc2WithErr(%s) returned error %v, expected %v", testName, err, expectedErr)
			}
		})

		if !strings.Contains(output, "[Timer] Timing "+testName) {
			t.Errorf("Log output missing initial timing message: %s", output)
		}
		if !strings.Contains(output, "[Timer] "+testName+" took") {
			t.Errorf("Log output missing final took message: %s", output)
		}
	})
}
