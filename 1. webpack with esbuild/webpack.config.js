// External modules
const path = require('path');

// Webpack plugins
const NodemonPlugin = require('nodemon-webpack-plugin');
const ForkTsCheckerPlugin = require('fork-ts-checker-webpack-plugin');

// Environment config
const isDevelopment = process.env.NODE_ENV !== 'production';
const mode = isDevelopment ? 'development' : 'production';

// Bundle config options
const BUNDLE = {
  entry: './index.ts',
  output: {
    filename: 'app.js',
    path: path.resolve(__dirname, 'dist')
  }
};

module.exports = {
  mode,
  target: 'node',
  entry: BUNDLE.entry,
  stats: 'errors-only',
  module: getLoaders(),
  plugins: getPlugins(),
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.json']
  },
  output: BUNDLE.output
};

/**
 * Loaders used by the application.
 */
function getLoaders() {
  const esbuild = {
    test: /\.(js|jsx|ts|tsx)?$/,
    loader: 'esbuild-loader',
    options: {
      loader: 'tsx',
      target: 'es2015'
    },
    exclude: /node_modules/
  };

  const loaders = {
    rules: [esbuild]
  };

  return loaders;
}

/**
 * Plugins
 */
function getPlugins() {
  const nodemon = new NodemonPlugin();
  const tsChecker = new ForkTsCheckerPlugin();

  return [tsChecker, nodemon];
}
