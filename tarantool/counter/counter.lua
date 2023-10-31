local log = require('log')

local counter = {
    start = function(self)
        -- create spaces and indexes
        box.once('init', function()
            box.schema.create_space('counters')
            box.space.counters:format({
                {name='user_id', type='string'},
                {name='id', type='string'},
                {name='value', type='integer'},
            })
            box.space.counters:create_index(
                    'primary', {type = 'tree', parts = {'user_id', 'id'}}
            )
        end)

        return true
    end,

    add = function(self, user_id, id, value)
        box.space.counters:upsert({user_id, id, value}, {{'+', 3, value}})
        return box.space.counters:get({user_id, id})[3]
    end,

    read = function(self, user_id, id)
        return box.space.counters:get({user_id, id})
    end,
}

return counter
