var Signature = Backbone.Model.extend({
    url: function() {
        var base = "sign";
        return base;
    }
});

var App = {
    Views: {},
    Controllers: {},
    init: function() {
        new App.Controllers.Signatures();
        Backbone.history.start();
    }
};

App.Views.Index = Backbone.View.extend({
    events: {
        "submit form": "sign"
    },

    initialize: function() {
        _.bindAll(this, 'render');
        this.model.bind('change', this.render);
        this.render();
    },

    render: function() {
        
        var compiled = _.template($('#index_template').html(), {model : this.model});
        $(this.el).html(compiled);
        $("#sandbar-form").html(this.el);
    },

    sign: function() {
        var self = this;
        this.model.save({
            Key: this.$('[name=Key]').val(),
            Content: this.$('[name=Content]').val()
        },
        {
            success: function(model, resp) {
                self.model = model;
                self.render();
                self.delegateEvents();
            },
            error: function() {
                alert('Failed to sign.');
            }
        });
        return false;
    }
});


App.Controllers.Signatures = Backbone.Controller.extend({
    routes: {
        "": "index"
    },

    index: function() {
        new App.Views.Index({
            model: new Signature()
        });
    }
});

$(function() {
    App.init()
});

window.JST