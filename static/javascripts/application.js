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
        var out = "<form class='sandbar-form'>";
        out += "<div class='field-label'>Content</div>";
        out += "<input class='textfield' type='text' name='Content' value='<%= model.get('Content') || ''%>'/>";
        out += "<div class='field-label'>Key</div>";
        out += "<input class='textfield' type='text' name='Key' value='<%= model.get('Key') || ''%>'/>";
        out += "<div class='buttons'><button class='sandbar-button'>Sign</button></div>";
        out += "<label>Signature: <%= model.get('Signature') || '' %></label</form>"
        var compiled = _.template(out);
        $(this.el).html(compiled({
            model: this.model
        }));
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