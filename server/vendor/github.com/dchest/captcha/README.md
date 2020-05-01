Package captcha
=====================

	import "github.com/dchest/captcha"

Package captcha implements generation and verification of image and audio
CAPTCHAs.

A captcha solution is the sequence of digits 0-9 with the defined length.
There are two captcha representations: image and audio.

An image representation is a PNG-encoded image with the solution printed on
it in such a way that makes it hard for computers to solve it using OCR.

An audio representation is a WAVE-encoded (8 kHz unsigned 8-bit) sound with the
spoken solution (currently in English, Russian, Chinese, and Japanese). To make
it hard for computers to solve audio captcha, the voice that pronounces numbers
has random speed and pitch, and there is a randomly generated background noise
mixed into the sound.

This package doesn't require external files or libraries to generate captcha
representations; it is self-contained.

To make captchas one-time, the package includes a memory storage that stores
captcha ids, their solutions, and expiration time. Used captchas are removed
from the store immediately after calling Verify or VerifyString, while
unused captchas (user loaded a page with captcha, but didn't submit the
form) are collected automatically after the predefined expiration time.
Developers can also provide custom store (for example, which saves captcha
ids and solutions in database) by implementing Store interface and
registering the object with SetCustomStore.

Captchas are created by calling New, which returns the captcha id.  Their
representations, though, are created on-the-fly by calling WriteImage or
WriteAudio functions. Created representations are not stored anywhere, but
subsequent calls to these functions with the same id will write the same
captcha solution. Reload function will create a new different solution for the
provided captcha, allowing users to "reload" captcha if they can't solve the
displayed one without reloading the whole page.  Verify and VerifyString are
used to verify that the given solution is the right one for the given captcha
id.

Server provides an http.Handler which can serve image and audio
representations of captchas automatically from the URL. It can also be used
to reload captchas.  Refer to Server function documentation for details, or
take a look at the example in "capexample" subdirectory.


Examples
--------

![Image](https://github.com/dchest/captcha/raw/master/capgen/example.png)

[Audio](https://github.com/dchest/captcha/raw/master/capgen/example.wav)


Constants
---------

``` go
const (
    // Default number of digits in captcha solution.
    DefaultLen = 6
    // The number of captchas created that triggers garbage collection used
    // by default store.
    CollectNum = 100
    // Expiration time of captchas used by default store.
    Expiration = 10 * time.Minute
)
```

``` go
const (
    // Standard width and height of a captcha image.
    StdWidth  = 240
    StdHeight = 80
)
```


Variables
---------

``` go
var (
    ErrNotFound = errors.New("captcha: id not found")
)
```



Functions
---------

### func New

	func New() string
	
New creates a new captcha with the standard length, saves it in the internal
storage and returns its id.

### func NewLen

	func NewLen(length int) (id string)
	
NewLen is just like New, but accepts length of a captcha solution as the
argument.

### func RandomDigits

	func RandomDigits(length int) (b []byte)
	
RandomDigits returns a byte slice of the given length containing
pseudorandom numbers in range 0-9. The slice can be used as a captcha
solution.

### func Reload

	func Reload(id string) bool
	
Reload generates and remembers new digits for the given captcha id.  This
function returns false if there is no captcha with the given id.

After calling this function, the image or audio presented to a user must be
refreshed to show the new captcha representation (WriteImage and WriteAudio
will write the new one).

### func Server

	func Server(imgWidth, imgHeight int) http.Handler
	
Server returns a handler that serves HTTP requests with image or
audio representations of captchas. Image dimensions are accepted as
arguments. The server decides which captcha to serve based on the last URL
path component: file name part must contain a captcha id, file extension —
its format (PNG or WAV).

For example, for file name "LBm5vMjHDtdUfaWYXiQX.png" it serves an image captcha
with id "LBm5vMjHDtdUfaWYXiQX", and for "LBm5vMjHDtdUfaWYXiQX.wav" it serves the
same captcha in audio format.

To serve a captcha as a downloadable file, the URL must be constructed in
such a way as if the file to serve is in the "download" subdirectory:
"/download/LBm5vMjHDtdUfaWYXiQX.wav".

To reload captcha (get a different solution for the same captcha id), append
"?reload=x" to URL, where x may be anything (for example, current time or a
random number to make browsers refetch an image instead of loading it from
cache).

By default, the Server serves audio in English language. To serve audio
captcha in one of the other supported languages, append "lang" value, for
example, "?lang=ru".

### func SetCustomStore

	func SetCustomStore(s Store)
	
SetCustomStore sets custom storage for captchas, replacing the default
memory store. This function must be called before generating any captchas.

### func Verify

	func Verify(id string, digits []byte) bool
	
Verify returns true if the given digits are the ones that were used to
create the given captcha id.

The function deletes the captcha with the given id from the internal
storage, so that the same captcha can't be verified anymore.

### func VerifyString

	func VerifyString(id string, digits string) bool
	
VerifyString is like Verify, but accepts a string of digits.  It removes
spaces and commas from the string, but any other characters, apart from
digits and listed above, will cause the function to return false.

### func WriteAudio

	func WriteAudio(w io.Writer, id string, lang string) error
	
WriteAudio writes WAV-encoded audio representation of the captcha with the
given id and the given language. If there are no sounds for the given
language, English is used.

### func WriteImage

	func WriteImage(w io.Writer, id string, width, height int) error
	
WriteImage writes PNG-encoded image representation of the captcha with the
given id. The image will have the given width and height.


Types
-----

``` go
type Audio struct {
    // contains unexported fields
}
```


### func NewAudio

	func NewAudio(id string, digits []byte, lang string) *Audio
	
NewAudio returns a new audio captcha with the given digits, where each digit
must be in range 0-9. Digits are pronounced in the given language. If there
are no sounds for the given language, English is used.

Possible values for lang are "en", "ja", "ru", "zh".

### func (*Audio) EncodedLen

	func (a *Audio) EncodedLen() int
	
EncodedLen returns the length of WAV-encoded audio captcha.

### func (*Audio) WriteTo

	func (a *Audio) WriteTo(w io.Writer) (n int64, err error)
	
WriteTo writes captcha audio in WAVE format into the given io.Writer, and
returns the number of bytes written and an error if any.

``` go
type Image struct {
    *image.Paletted
    // contains unexported fields
}
```


### func NewImage

	func NewImage(id string, digits []byte, width, height int) *Image
	
NewImage returns a new captcha image of the given width and height with the
given digits, where each digit must be in range 0-9.

### func (*Image) WriteTo

	func (m *Image) WriteTo(w io.Writer) (int64, error)
	
WriteTo writes captcha image in PNG format into the given writer.

``` go
type Store interface {
    // Set sets the digits for the captcha id.
    Set(id string, digits []byte)

    // Get returns stored digits for the captcha id. Clear indicates
    // whether the captcha must be deleted from the store.
    Get(id string, clear bool) (digits []byte)
}
```

An object implementing Store interface can be registered with SetCustomStore
function to handle storage and retrieval of captcha ids and solutions for
them, replacing the default memory store.

It is the responsibility of an object to delete expired and used captchas
when necessary (for example, the default memory store collects them in Set
method after the certain amount of captchas has been stored.)

### func NewMemoryStore

	func NewMemoryStore(collectNum int, expiration time.Duration) Store
	
NewMemoryStore returns a new standard memory store for captchas with the
given collection threshold and expiration time in seconds. The returned
store must be registered with SetCustomStore to replace the default one.
