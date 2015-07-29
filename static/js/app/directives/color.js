angular.module('zfh').directive('myColor', function(){
    return {
        restrict : "A",
        scope : {
            color : "=myColor"
        },
        link : function(scope, element, atts){
            element.ColorPicker({
                onShow : function (colpkr) {
                    $(colpkr).fadeIn(500);
                    return false;
                },
                onHide : function (colpkr) {
                    $(colpkr).fadeOut(500);
                    return false;
                },
                onChange : $.proxy(function (hsb, hex, rgb) {
                    var color = '#' + hex;
                    this.css('backgroundColor', color).next().val(color);
                    scope.$apply(function(){
                        scope.color = color;
                    });
                }, element)
            });
        }
    }
});