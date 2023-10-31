box.cfg({listen="0.0.0.0:3301"})
box.schema.user.create('storage', {password='passw0rd', if_not_exists=true})
box.schema.user.grant('storage', 'super', nil, nil, {if_not_exists=true})
require('msgpack').cfg{encode_invalid_as_nil = true}

counter = require('counter')
counter:start()