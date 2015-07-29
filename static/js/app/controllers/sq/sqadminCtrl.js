angular.module('zfh').controller('ctrl.sqadmin', ['$scope', '$http', '$timeout', function($scope, $http, $timeout){
    var vm = $scope.vm = {};

    var statusMap = {
        '0' : '申请中',
        '1' : '成功',
        '2' : '被拒绝'
    };

    $.extend(vm, {
        sqs : {},

        init : function(reset){
            !reset && $.each(Sqs, function(_, sq){
                sq['statusString'] = statusMap[sq['Status']];
                vm.sqs[sq.Name] = sq;
            });
        },

        ajaxing : "false",
        ajaxdo : function(status){
            vm.ajaxing = status;
            if(status != "true" || status != "false") {
                $timeout(function(){
                    vm.ajaxing = "false";
                }, 0);
            }
        },

        status : function(name, status){
            vm.doAjax('/sq/status', {name : name, status : status}, function(data){
                $.extend(vm.sqs[name], {Status : status, statusString : statusMap[status]});
            });
            return false;
        },

        doAjax : function(url, post, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(post),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    cb && cb(data);
                    vm.ajaxdo("false");
                }
            });
        }
    });

    vm.init();
}]);