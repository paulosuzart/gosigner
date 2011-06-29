var signView = {
    Key : ko.observable(''),
    Signature : ko.observable(''),
    Content : ko.observable(''),
    sign : function () {
       $.post(  'sign', 
                ko.toJSON(signView), 
                function(data){
                   ko.mapping.fromJS(data, {}, signView); 
                   $("#sign").css('visibility', 'visible');
                });
    },
};
/**
var App = { signView : signView, 
            aboutView : aboutView,
            about : false,
            display : function(about) {
              return about? "index_template" : "about_template";  
            },
            about : function(){
                App.about = true;
            }
        }
*/
$(function(){
    ko.applyBindings(signView, document.getElementById('sandbar-form'));  
});
