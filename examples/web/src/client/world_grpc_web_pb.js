/**
 * @fileoverview gRPC-Web generated client stub for world
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.world = require('./world_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.world.HelloClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.world.HelloPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!proto.world.HelloClient} The delegate callback based client
   */
  this.delegateClient_ = new proto.world.HelloClient(
      hostname, credentials, options);

};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.world.HelloRequest,
 *   !proto.world.HelloReply>}
 */
const methodInfo_Hello_Visit = new grpc.web.AbstractClientBase.MethodInfo(
  proto.world.HelloReply,
  /** @param {!proto.world.HelloRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.world.HelloReply.deserializeBinary
);


/**
 * @param {!proto.world.HelloRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.world.HelloReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.world.HelloReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.world.HelloClient.prototype.visit =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/world.Hello/Visit',
      request,
      metadata,
      methodInfo_Hello_Visit,
      callback);
};


/**
 * @param {!proto.world.HelloRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.world.HelloReply>}
 *     The XHR Node Readable Stream
 */
proto.world.HelloPromiseClient.prototype.visit =
    function(request, metadata) {
  return new Promise((resolve, reject) => {
    this.delegateClient_.visit(
      request, metadata, (error, response) => {
        error ? reject(error) : resolve(response);
      });
  });
};


module.exports = proto.world;

