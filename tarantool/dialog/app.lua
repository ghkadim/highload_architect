local dialog = require('dialog')
local json = require('json')
box.cfg{}

local function dialog_add(req)
    msg_id = dialog:add(req:json())
    return {
        status = 200,
        body = json.encode({message_id=msg_id})
    }
end

local function dialog_list(req)
    from = req:param('from')
    to = req:param('to')
    if from == nil or to == nil then
        error({code=400})
    end
    result = dialog:list(from, to)
    return {
        status = 200,
        body = json.encode(result)
    }
end

local function dialog_read_message(req)
    message_id = tonumber(req:param('message_id'))
    read_by_user = req:param('user_id')
    if message_id == nil or read_by_user == nil then
        error({code=400})
    end
    dialog:read(message_id, read_by_user)
    return {
        status = 200,
    }
end

local function err_middleware(f)
    return function(req)
        status, res = pcall(f, req)
        if status then
            return res
        else
            err = res
            if type(err) == "table" and err.code ~= nil then
                return {
                    status = err.code,
                    body = json.encode({message=err.message})
                }
            else
                return {
                    status = 500,
                    body = json.encode({message=res})
                }
            end
        end
    end
end

dialog:start()

local server = require('http.server').new(nil, 3380, {charset = "utf8"})
server:route({path = '/dialog', method = 'GET'}, err_middleware(dialog_list))
server:route({path = '/dialog', method = 'POST'}, err_middleware(dialog_add))
server:route({path = '/dialog/message/read', method = 'PUT'}, err_middleware(dialog_read_message))
server:start()