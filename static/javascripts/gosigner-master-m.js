var signView = {
    Key : ko.observable(''),
    Signature : ko.observable(''),
    Content : ko.observable(''),
    sign : function () {
       $.post(  'sign', 
                ko.toJSON(signView), 
                function(data){
                   ko.mapping.fromJS(data, {}, signView);
                   //$("#sign").css('visibility', 'visible');
                });
    }
};

$(function(){
    ko.applyBindings(signView);
});
