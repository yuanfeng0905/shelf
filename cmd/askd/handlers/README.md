

# handlers
`import "github.com/coralproject/shelf/cmd/askd/handlers"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package handlers contains the handler logic for processing requests.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)


#### <a name="pkg-files">Package files</a>
[form.go](/src/github.com/coralproject/shelf/cmd/askd/handlers/form.go) [form_gallery.go](/src/github.com/coralproject/shelf/cmd/askd/handlers/form_gallery.go) [form_submission.go](/src/github.com/coralproject/shelf/cmd/askd/handlers/form_submission.go) [version.go](/src/github.com/coralproject/shelf/cmd/askd/handlers/version.go) 



## <a name="pkg-variables">Variables</a>
``` go
var ErrInvalidCaptcha = errors.New("captcha invalid")
```
ErrInvalidCaptcha is returned when a captcha is required for a form but it
is not valid on the request.

``` go
var Form formHandle
```
Form fronts the access to the form service functionality.

``` go
var FormGallery formGalleryHandle
```
FormGallery fronts the access to the form service functionality.

``` go
var FormSubmission formSubmissionHandle
```
FormSubmission fronts the access to the form service functionality.

``` go
var Version verHandle
```
Version fronts the access to the ver service functionality.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)