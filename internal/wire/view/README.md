

# view
`import "github.com/coralproject/shelf/internal/wire/view"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Delete(context interface{}, db *db.DB, name string) error](#Delete)
* [func GetAll(context interface{}, db *db.DB) ([]View, error)](#GetAll)
* [func Upsert(context interface{}, db *db.DB, view *View) error](#Upsert)
* [type Path](#Path)
  * [func (slice Path) Len() int](#Path.Len)
  * [func (slice Path) Less(i, j int) bool](#Path.Less)
  * [func (slice Path) Swap(i, j int)](#Path.Swap)
* [type PathSegment](#PathSegment)
  * [func (ps *PathSegment) Validate() error](#PathSegment.Validate)
* [type View](#View)
  * [func GetByName(context interface{}, db *db.DB, name string) (*View, error)](#GetByName)
  * [func (v *View) Validate() error](#View.Validate)


#### <a name="pkg-files">Package files</a>
[model.go](/src/github.com/coralproject/shelf/internal/wire/view/model.go) [view.go](/src/github.com/coralproject/shelf/internal/wire/view/view.go) 


## <a name="pkg-constants">Constants</a>
``` go
const Collection = "views"
```
Collection is the Mongo collection containing view metadata.


## <a name="pkg-variables">Variables</a>
``` go
var ErrNotFound = errors.New("Set Not found")
```
ErrNotFound is an error variable thrown when no results are returned from a Mongo query.



## <a name="Delete">func</a> [Delete](/src/target/view.go?s=2458:2520#L81)
``` go
func Delete(context interface{}, db *db.DB, name string) error
```
Delete removes a view from from Mongo.



## <a name="GetAll">func</a> [GetAll](/src/target/view.go?s=1218:1277#L36)
``` go
func GetAll(context interface{}, db *db.DB) ([]View, error)
```
GetAll retrieves the current views from Mongo.



## <a name="Upsert">func</a> [Upsert](/src/target/view.go?s=487:548#L10)
``` go
func Upsert(context interface{}, db *db.DB, view *View) error
```
Upsert upserts a view to the collection of currently utilized views.




## <a name="Path">type</a> [Path](/src/target/model.go?s=875:898#L20)
``` go
type Path []PathSegment
```
Path is a slice of PathSegment.










### <a name="Path.Len">func</a> (Path) [Len](/src/target/model.go?s=1127:1154#L31)
``` go
func (slice Path) Len() int
```
Len is required to sort a slice of PathSegment.




### <a name="Path.Less">func</a> (Path) [Less](/src/target/model.go?s=1231:1268#L36)
``` go
func (slice Path) Less(i, j int) bool
```
Less is required to sort a slice of PathSegment.




### <a name="Path.Swap">func</a> (Path) [Swap](/src/target/model.go?s=1366:1398#L41)
``` go
func (slice Path) Swap(i, j int)
```
Swap is required to sort a slice of PathSegment.




## <a name="PathSegment">type</a> [PathSegment](/src/target/model.go?s=517:838#L12)
``` go
type PathSegment struct {
    Level     int    `bson:"level" json:"level" validate:"required,min=1"`
    Direction string `bson:"direction" json:"direction" validate:"required,min=2"`
    Predicate string `bson:"predicate" json:"predicate" validate:"required,min=1"`
    Tag       string `bson:"tag,omitempty" json:"tag,omitempty"`
}
```
PathSegment contains metadata about a segment of a path,
which path partially defines a View.










### <a name="PathSegment.Validate">func</a> (\*PathSegment) [Validate](/src/target/model.go?s=958:997#L23)
``` go
func (ps *PathSegment) Validate() error
```
Validate checks the PathSegment value for consistency.




## <a name="View">type</a> [View](/src/target/model.go?s=1485:1813#L46)
``` go
type View struct {
    Name       string `bson:"name" json:"name" validate:"required,min=3"`
    Collection string `bson:"collection" json:"collection" validate:"required,min=2"`
    StartType  string `bson:"start_type" json:"start_type" validate:"required,min=3"`
    Path       Path   `bson:"path" json:"path" validate:"required,min=1"`
}
```
View contains metadata about a view.







### <a name="GetByName">func</a> [GetByName](/src/target/view.go?s=1801:1875#L58)
``` go
func GetByName(context interface{}, db *db.DB, name string) (*View, error)
```
GetByName retrieves a view by name from Mongo.





### <a name="View.Validate">func</a> (\*View) [Validate](/src/target/model.go?s=1866:1897#L54)
``` go
func (v *View) Validate() error
```
Validate checks the View value for consistency.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
