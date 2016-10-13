var gulp = require('gulp')
var plugins = require('gulp-load-plugins')()
var browserify = require('browserify')
var watchify = require('watchify')
var source = require('vinyl-source-stream')
var buffer = require('vinyl-buffer')
var babelify = require('babelify')
var pkg = require(__dirname + '/package.json')
var path = require('path')
var es = require('event-stream')
var sass = require('gulp-sass')


gulp.task('build', function() {
  return compile(false)
})

gulp.task('build-release', function() {
  return compile(true)
})

gulp.task('watchify', function() {
  return compile(false, true)
})

gulp.task('sass', function() {
  return sassCompile('expanded')
})

gulp.task('sass-release', function() {
  return sassCompile('compressed')
})

gulp.task('watch', ['watchify', 'sass'], function() {
  gulp.watch(['./frontend/stylesheets/**/*.scss'], ['sass'])
})

function compile(isUglify, isWatch) {
  var entries = pkg.browserify.entryScripts

  var tasks = entries.map( function(entry) {
    var extname = path.extname(entry)
    var output = path.basename(entry, extname) + '.js'
    var bundler = null

    bundler = getBundler(entry, isWatch)
    function bundle() {
      return bundler.bundle()
        .on('error', plugins.util.log)
        .on('error', handleError)
        .pipe(source(output))
        .pipe(buffer())
        .pipe(plugins.sourcemaps.init({ loadMaps: true }))
        .pipe(plugins.duration("compiled :" + entry))
        .pipe(plugins.if(isUglify, plugins.uglify()))
        .pipe(plugins.sourcemaps.write('.'))
        .pipe(gulp.dest("public/assets/javascripts/"))
    }
    bundler.on('update', bundle)
    return bundle()
  })
  es.merge.apply(null, tasks)
}

function getBundler(entry, isWatch) {
  var bundler = null
  var options = { debug: true }
  if (isWatch) {
    options.cache = {}
    options.packageCache = {}
    options.fullPaths = true
  }
  bundler = browserify(entry, options)

  if (isWatch) {
    bundler = watchify(bundler, {poll: true})
  }

  bundler.transform(babelify, {presets: ["es2015", "react"]})
  return bundler
}

function handleError() {
  this.emit('end')
}

function sassCompile(sassStyle) {
  return gulp.src("frontend/stylesheets/*.scss")
    .pipe(plugins.plumber())
    .pipe(plugins.sourcemaps.init())
    .pipe(sass({
      style: sassStyle
    }))
    .pipe(plugins.sourcemaps.write('.'))
    .pipe(gulp.dest("public/assets/stylesheets"))
}
