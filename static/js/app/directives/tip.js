angular.module('zfh').directive('myTip', function(){
    return {
        restrict : "A",
        link : function(scope, element, atts){
            scope.$watch(atts.myTip, function(value){
                if(value == "true" || value == "false"){
                    element.triggerHandler('load', [value]);
                }else{
                    element.triggerHandler('tip', [value]);
                }
            });

            element.on({
                load : function(event, status){
                    $(this)[status == "true" ? 'addClass' : 'removeClass']('ajax-loading');
                },

                tip : function(event, string){
                    var $this = $(this);
                    $this.triggerHandler('load', ["false"]);
                    var pp = $this.offset();
                    var pw = $this.outerWidth();
                    var ph = $this.outerHeight();
                    var tipBox = $('<div/>').appendTo('body').css({
                        position : "absolute",
                        left : pp.left + pw/2 + "px",
                        top : pp.top + "px",
                        zIndex : 1000001,
                        width : 0,
                        height : 0,
                        'text-align' : 'center'
                    });
                    var innerTip = $('<div/>').text(string).appendTo(tipBox).css({
                        position : "relative",
                        top : 0,
                        color : 'red',
                        'white-space' : 'nowrap'
                    });
                    setTimeout(function(){
                        innerTip.css('left', '-' + innerTip.outerWidth() / 2 + 'px');
                    }, 0);
                    innerTip.animate({top : '-25'}, 250, function(){
                        var $this = $(this);
                        setTimeout(function(){
                            $this.animate({top : 0, opacity : 0}, 500, function(){
                                $(this).parent().remove();
                            });
                        }, 1000);
                    });
                }
            });
        }
    }
});