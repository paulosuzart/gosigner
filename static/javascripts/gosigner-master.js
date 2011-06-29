
var signView = {
    Key : ko.observable(''),
    Signature : ko.observable(''),
    Content : ko.observable(''),
    sign : function () {
       $.post(  'sign', 
                ko.toJSON(this), 
                function(data){
                   ko.mapping.fromJS(data, {}, signView); 
                   $("#sign").css('visibility', 'visible');
                });
    },
};

$(function(){
    ko.applyBindings(signView, document.getElementById('sandbar-form'));  
});
