const {
  HelloRequest
} = require('./client/world_pb')
const {
  HelloClient
} = require('./client/world_grpc_web_pb')

const grpc = {};
grpc.web = require('grpc-web')

const grpcurl = 'http://localhost:18080'
var helloService = new HelloClient(grpcurl, null, null)

console.log(helloService)

var req = new HelloRequest("from web")

var metadata = {}
helloService.visit(req, metadata, (err, res) => {
  console.log(err, res)
})