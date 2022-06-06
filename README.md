# tower-colour-api
Change the colour of a Dell G5 tower

A cross platform CLI and web API for updating the LED colours of a Dell G5 tower.

Example at [chez.kindrobot.ca/puter.html](https://chez.kindrobot.ca/puter.html).

## Installing
Precompiled binaries are available on the releases page.

### Building

This can be built from source with a modern [golang] compiler. On Windows, you will need to have [gcc for Windows] installed.

```bash
  go build
```

## Usage

```bash
./tower-colour-api RED GREEN BLUE
```

Updates the colour of the tower. `RED`, `GREEN`, and `BLUE` must be integers between 0 and 255.

```bash
./tower-colour-api LOG_FILE_PATH
```

Starts a server on port 3010 to listen for HTTP requests to update the tower colour.

## API

```
POST :3010/
{
  red: 255,
  green: 255,
  blue: 255,
  message: "a nice note"
}
```

Updates the colour of the tower. All fields are optional. `red`, `green`, and `blue` must be integers between 0 and 255 and default to 0. `message` is logged to the log fle along with the colour and time of request.

```
GET :3010/
```

Shows the last colour update in the form:

```
{
  red: 255,
  green: 255,
  blue: 255,
  message: "a nice note",
  created_at: "2022-06-06Z15:42.12345"
}
```

## Other files

| File        | Purpose                                                                  |
|-------------|--------------------------------------------------------------------------|
| puter.html  | A simple HTML frontend for the colour change API                         |
| motion.conf | Config for a webcam using [motion] to update the picture on the frontend |

## Credits
All of the HID calls were taken from [T-Troll/alienfx-tools](https://github.com/T-Troll/alienfx-tools).  

[golang]: https://go.dev/doc/install
[gcc for Windows]: https://www.mingw-w64.org/
[motion]: https://motion-project.github.io/
