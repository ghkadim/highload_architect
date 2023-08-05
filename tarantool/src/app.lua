local dialog = require('dialog')
local json = require('json')
box.cfg{}

local function dialog_add(req)
    dialog:add(req:json())
    return {status = 200}
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
server:start()