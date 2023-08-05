local log = require('log')
local avro = require('avro_schema')

local schema = {
    dialog = {
        type='record',
        name='dialog_schema',
        fields={
            {name='id', type='long*'},
            {name='from', type='string'},
            {name='to', type='string'},
            {name='text', type='string'}
        }
    }
}

local dialog = {
    dialog_model = {},

    start = function(self)
        -- create spaces and indexes
        box.once('init', function()
            box.schema.sequence.create('dialog_id')
            box.schema.create_space('dialogs')
            box.space.dialogs:format({
                {name='id', type='unsigned'},
                {name='from', type='string'},
                {name='to', type='string'},
                {name='text', type='string'},
                {name='search_key', type='string'}
            })
            box.space.dialogs:create_index(
                    'primary', {type = 'tree', parts = {'id'}, sequence='dialog_id'}
            )
            box.space.dialogs:create_index(
                    'search_key', {type = 'tree', parts = {'search_key', 'id'}}
            )
        end)

        -- create models
        local ok, dialog = avro.create(schema.dialog)
        if ok then
            -- compile models
            local ok, compiled_dialog = avro.compile(dialog)
            if ok then
                self.dialog_model = compiled_dialog
                log.info('Started')
                return true
            else
                log.error('Schema compilation failed')
            end
        else
            log.info('Schema creation failed')
        end
        return false
    end,

    -- return dialog list
    list = function(self, user1, user2)
        local result = {}
        local sk = search_key(user1, user2)
        for _, tuple in box.space.dialogs.index.search_key:pairs(sk, {iterator = box.index.LE}) do
            if (tuple[5] ~= sk) then break end
            local _, dialog = self.dialog_model.unflatten({tuple[1], tuple[2], tuple[3], tuple[4]})
            table.insert(result, dialog)
        end
        return result
    end,

    -- add dialog
    add = function(self, dialog)
        local ok, tuple = self.dialog_model.flatten(dialog)
        if not ok then
            error({code = 400, message = tuple})
        end
        tuple[5] = search_key(dialog.from, dialog.to)
        box.space.dialogs:insert(tuple)
    end,
}

function search_key(from_user, to_user)
    if from_user < to_user then
        return from_user .. ':' .. to_user
    else
        return to_user .. ':' .. from_user
    end
end

return dialog