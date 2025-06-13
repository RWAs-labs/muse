# Wait

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-wait)](https://github.com/nhatthm/go-wait/releases/latest)
[![Build Status](https://github.com/nhatthm/go-wait/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-wait/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-wait/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-wait)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/wait)](https://goreportcard.com/report/go.nhat.io/wait)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/wait)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

A simple library to wait for something.

## Prerequisites

- `Go >= 1.17`

## Install

```bash
go get go.nhat.io/wait
```

## Usage

```go
package main

import (
	"context"
	"time"

	"go.nhat.io/wait"
)

func execute(ctx context.Context) error {
	if err := wait.ForDuration(time.Minute).Wait(ctx); err != nil {
        return err
	}

	// Do something here.

	return nil
}
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
