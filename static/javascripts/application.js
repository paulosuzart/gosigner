var Signature = Backbone.Model.extend({
    url: function() {
        var base = "api/sign";
        return base;
    },
    validate: function(attrs){
        if (_(attrs.Content).trim().length == 0) {
                return "You should put some conent to be signed :)."
        }
        if ((attrs.Key).trim().length == 0) {
                return "Where is the Key to sign?"
        }
    }
});

var Key = Backbone.Model.extend({
        url: function(){
                var base = "api/key";
                if(this.isNew()) return base;
                return base + (base.charAt(base.length - 1) == '/' ? '' : '/') + this.get('alias');
        }
});

var Keys = Backbone.Collection.extend({
        model: Key,
        url : "api/keys"
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
        "submit form" : "sign",
        "click [name=Content]" : "clear" 
    },

    clear: function(){
      $("#sign").css('visibility','hidden');      
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
                $("#sign").css('visibility', 'visible');
                self.delegateEvents();
            },
            error: function(self, error) {
                if (error) {
                        alert(error);
                } else {
                        alert('Failed to sign.');
                }
            }
        });
        return false;
    },
});


App.Controllers.Signatures = Backbone.Controller.extend({
    routes: {
        "": "index",
        "/about" : "about"
    },

    index: function() {
        new App.Views.Index({
            model: new Signature()
        });
    },

    about : function() {
        $.get('/api/ver', function(data){
                var compiled = _.template($('#about_template').html(), {version : data.Version});
                $("#sandbar-form").html(compiled);
                Backbone.history.saveLocation('/about');
        }, "json");
    }
});

$(function() {
    App.init()
});
