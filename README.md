# JS Unpacker
This program attempts to harvest as much information as possible from javascript source maps.
Works with both, local files and urls. See [output](#output) for additional detail. 

```bash
go install github.com/webklex/juck
```

## Usage
```bash
Usage of juck:
  --file      string    Target sourcemap file path
  --file-list string    File path of a file containing a list of target source map file paths
  --url       string    Target sourcemap url
  --url-list  string    File path of a file containing a list of target source map urls
  --force               Force to download and overwrite local sourcemap
  --delay     duration  Delay between two requests. Only applies if --url-list is used
  --output    string    Directory to output from sourcemap to (default "./output")
  --combined            Combine all source files into one
  --disable-ssl         Don't verify the site's SSL certificate
  --no-color            Disable color output
  --version             Show version and exit
  --dangerously-write-paths  Write full paths. WARNING: Be careful here, you are pulling directories from an untrusted source
```

Analyze a single local file:
```bash
./juck --file ./source.js.map
```

Analyze a single url:
```bash
./juck --url https://example.com/assets/some_file.js
```
> Note: you don't have to apply a .map - it gets added automatically if it is missing.


Analyze a file containing many urls and delay each request by 3 seconds:
```bash
cat ./url_list.txt
https://example.com/assets/js/some_file.js
https://example.com/some_other_file.js
...
```
```bash
./juck --url-list ./url_list.txt --delay 3s
```

Analyze piped stdin:
```bash
echo "https://example.com/assets/js/some_file.js" | ./juck
```
..or:
```bash
echo "./source.js.map" | ./juck
```

## Output
By default, the output is stored in a folder called `output` placed within your current working directory.
The output folder contains the following folders and files after the program has run:
- `combined` - all combined files (only if `--combined` is active)
- `sourcemaps` - all downloaded source maps
- `sources` - all recovered sources
- `node_modules.txt` - a list of all directly discovered node modules
- `dependencies.txt` - a list of all additional dependencies based on the latest version registered on [www.npmjs.com](https://www.npmjs.com/)


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