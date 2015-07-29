angular.module('zfh').controller('ctrl.tskl', ['$scope', '$http', '$timeout', '$interval', "$sce", "myTime", function($scope, $http, $timeout, $interval, $sce, myTime){
    var vm = $scope.vm = {};
    $.extend(vm, {
        tskl : {},
        state : [],
        jsDemo : 0,

        pause : false,

        step : 0,

        zcs : [],
        backjson : false,
        tsChange : 0,

        init : function(reset){
            if(!reset) {
                vm.tskl = TSKL /*{
                    back : JSON.stringify({"a" : '1', "b" : '2'}),
                    berr : "222"
                }*/;
                var step = vm._stepStore();
                if(step > 0){
                    vm._step(step);
                }
                vm.zcs = vm._store(true);
            }
            vm.ajax();
        },

        ajax : function(){
            $http({
                method : "GET",
                url : '/tskl/ajax',
                responseType : "json",
                params : {"hash" : vm.tskl.hash}
            }).success(function(data, status){
                if(data.error){
                    console.log(data.error);
                } else {
                    if(!data['hash']){
                        vm.pause = true;
                    }
                    $.extend(vm.tskl, data);
                }
            }).finally(function(){
                if(vm.pause){
                    return;
                }
                $timeout(function(){vm.ajax();}, 100);
            });
        },

        send : function(){
            if(!vm.tskl.send || vm.step > 0){
                return;
            }
            $http({
                method : "POST",
                url : '/tskl/post',
                responseType : "json",
                data : $.param(vm.tskl)
            }).success(function(data, status){
                if(data.error){
                    console.log(data.error);
                } else {
                    vm.pause = false;
                    vm.tskl.hash = data['hash'];
                    vm.tskl.back = vm.tskl.berr = '';
                    vm.ajax();
                    vm._step();
                }
            });
        },

        _step : function(ss){
            vm.step = ss || STEP;
            var timer = $interval(function(){
                vm.step--;
                if(vm.step <= 0){
                    vm.step = 0;
                    $interval.cancel(timer);
                    timer = null;
                }
                vm._stepStore(vm.step);
            }, 1000);
        },

        _stepStore : function(val){
            var key = 'tskl:step';
            if(val == undefined){
                return localStorage.getItem(key);
            }
            localStorage.setItem(key, val);
        },

        option : function(which){
            vm[which + 'Demo'] = vm[which + 'Demo'] ? 0 : 1;
        },

        zcbtn : function(){
            if(!vm.tskl.back) return;
            vm.zcs.push({back:vm.tskl.back});
            vm.zcs.length > 3 && (vm.zcs.length = 3);
            vm._store();
        },

        zcrm : function(i){
            vm.zcs.splice(i, 1);
            vm._store();
        },

        _store : function(bool){
            if(!bool) {
                var zcs = [];
                $.each(vm.zcs, function (i, v) {
                    delete(v['$$hashKey']);
                    zcs.push(v);
                });
                localStorage.setItem('zcs', JSON.stringify(zcs));
            } else {
                return JSON.parse(localStorage.getItem('zcs')) || [];
            }
        },

        jsonbtn : function(){
            if(vm.tskl.back) {
                //var jsonHtml = vm.jsonToHtml(vm.tskl.back, true, true);
                //vm.backjson = $sce.trustAsHtml(jsonHtml);
                vm.tsChange++;
                vm.backjson = !vm.backjson;
            }
        },

        jsonToHtml : function(json, collapsed, nl2br){
            var format = new JSONFormatter({
                collapsed : collapsed,
                nl2br: nl2br
            });
            return format.jsonToHTML(json);
        }
    });
    vm.init();
}]);