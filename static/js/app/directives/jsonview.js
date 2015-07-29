angular.module('zfh').directive('myJsonview', function(){
    return {
        restrict : "A",
        link : function(scope, element, atts){
            var tsChange = 0;
            scope.$watch('vm.tsChange', function(newV, oldV){
                if(newV != tsChange) {
                    setTimeout(function() {
                        $(element).find('.json-view').not('.json-view-set').each(function () {
                            try {
                                $(this).addClass('json-view-set').JSONView($(this).text(), {
                                    collapsed: false,
                                    nl2br: true
                                });
                            }catch(e){}
                        });
                    }, 0);
                }
            });
        }
    }
});