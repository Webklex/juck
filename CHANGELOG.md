# Changelog

All notable changes to `webklex/juck` will be documented in this file.

Updates should follow the [Keep a CHANGELOG](http://keepachangelog.com/) principles.


## [UNRELEASED]
### Fixed
- NaN

### Added
- NaN

### Breaking changes
- NaN


## [1.2.0] - 2022-09-10
### Added
- Verify npm package names
- Fetch **all** dependencies
- Accept variable targets (filename or url) from stdin
- Log level support added to suppress noise
- Also looks for css source maps

### Breaking changes
- Output folder content and structure has changed. All discovered sources are now inside a `sources` folder within the given output folder.


## [1.1.0] - 2022-09-09
### Fixed
- Add suffix to indirect resource references
- Ignore empty urls

### Added
- Delay downloads (e.g.: --delay 3s)
- Download sourcemaps to `{output}/sourcemaps`
- Download log added
- Local sourcemap cache added
- Optionally only use cached files
- Search for possible node module dependencies

### Breaking changes
- NaN


## [1.0.0] - 2022-09-06
Initial release
