# JS Unpacker
Recover uncompiled TypeScript sources, JSX, and more from Webpack sourcemaps. This program was heavily inspired by 
[rarecoil/unwebpack-sourcemap](https://github.com/rarecoil/unwebpack-sourcemap).

```bash
go install github.com/webklex/juck
```

## Usage
```bash
Usage of juck:
  --file      string  Target sourcemap file path
  --file-list string  File path of a file containing a list of target source map file paths
  --url       string  Target sourcemap url
  --url-list  string  File path of a file containing a list of target source map urls
  --output    string  Directory to output from sourcemap to (default "./output")
  --combined          Combine all source files into one
  --disable-ssl       Don't verify the site's SSL certificate
  --no-color          Disable color output
  --version           Show version and exit
  --dangerously-write-paths  Write full paths. WARNING: Be careful here, you are pulling directories from an untrusted source
```

Example:
```bash
./juck --file ./source.js.map
```

## Build
```bash
git clone https://github.com/webklex/juck
cd juck
go build
```
..or:
```bash
git clone https://github.com/webklex/juck
cd juck
./build.sh
```

## Security
If you discover any security related issues, please email github@webklex.com instead of using the issue tracker.

## Credits
- [Webklex][link-author]
- [rarecoil/unwebpack-sourcemap](https://github.com/rarecoil/unwebpack-sourcemap)
- [All Contributors][link-contributors]

## License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.

[link-author]: https://github.com/webklex
[link-contributors]: https://github.com/webklex/juck/graphs/contributors