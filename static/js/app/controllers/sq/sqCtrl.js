angular.module('zfh').controller('ctrl.sq', ['$scope', '$http', '$timeout', function($scope, $http, $timeout){
    var vm = $scope.vm = {};

    var statusMap = {
        '0' : '申请中',
        '1' : '成功',
        '2' : '被拒绝'
    };

    $.extend(vm, {
        sqs : [],
        ccsq : {},

        num : 0,
        maxNum : 1,
        init : function(reset){
            vm.ccsq = {
                name : "",
                loginName : "",
                loginPass : ""
            };

            !reset && $.each(Sqs, function(_, sq){
                sq['statusString'] = statusMap[sq['status']];
                vm.sqs.push(sq);
            });
            var num = parseInt(vm.store()) || 0;
            vm.num = num > vm.maxNum ? -1 : num;
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

        checkError : function(){
            var error = false;
            $.each(vm.ccsq, function(k, v){
                if($.trim(v).length == 0){
                    error = vm.config["null"];
                    return false;
                }
                if(k == 'loginName' && $.trim(v).length < 3){
                    error = vm.config["no"];
                    return false;
                }
            });
            return error;
        },

        doSQ : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/sq/post', function(data){
                vm.sqs.push($.extend({}, vm.ccsq, {statusString : statusMap['0']}));
                vm.store(++vm.num);
            });
            return false;
        },

        doAjax : function(url, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(vm.ccsq),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    cb && cb(data);
                    vm.init(true);
                    vm.ajaxdo("false");
                }
            });
        },

        store : function(val){
            var key = 'sq:num';
            if(val == undefined){
                return localStorage.getItem(key);
            }
            localStorage.setItem(key, val);
        },

        config : {
            "null" : "每项都需要填写",
            "no" : "别闹了，登录还是用你的姓名拼音吧，不能少于3个字母数字"
        }
    });

    vm.init();
}]);