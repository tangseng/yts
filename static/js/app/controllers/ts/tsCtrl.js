angular.module('zfh').controller('ctrl.ts', ['$scope', '$http', '$timeout', "myTime", function($scope, $http, $timeout, myTime){
    var vm = $scope.vm = {};

    var state = function(bool){
        return {
            v : bool,
            n : bool ? '简单' : '详细'
        }
    };
    $.extend(vm, {
        tss : [],
        hash : "",
        tsChange : 0,

        state : [],

        phpDemo : 0,
        jsDemo : 0,

        init : function(reset){
            if(!reset) {
                vm.tss = [];
                TS && $.each(TS, function (_, ts) {
                    ts.tstime = myTime.timeToStr(ts.tstime, true);
                    var data = JSON.parse(decodeURIComponent(ts.data));
                    var tmp = [];
                    if(ts.type == 'php') {
                        angular.forEach(data, function (v, k) {
                            var debug = v.debug;
                            var realData = debug[0]['args'];
                            if (debug.length == 1) {
                                debug = debug[0];
                            }

                            this.push({
                                time: myTime.timeToStr(v.time, 1),
                                info: JSON.stringify(debug),
                                realInfo: realData
                            });
                        }, tmp);
                    }else{
                        angular.forEach(data, function (v, k) {
                            var debug = v.debug;
                            if (debug.length == 1) {
                                debug = debug[0];
                            }
                            this.push({
                                time: myTime.timeToStr(v.time, 1),
                                info: JSON.stringify(debug)
                            });
                        }, tmp);
                    }
                    ts.data = tmp;
                    vm.tss.push(ts);
                    vm.state.push(state(0));
                });
                vm.tsChange++;
                vm.hash = Hash;
            }
        },

        ajax : function(){
            $http({
                method : "GET",
                url : '/ts/ajax',
                responseType : "json",
                params : {"hash" : vm.hash}
            }).success(function(data, status){
                if(data.error){
                    console.log(data.error);
                } else {
                    if(data['hash'] == vm.hash){
                        return;
                    }
                    TS = data["ts"];
                    Hash = data["hash"];
                    vm.init();
                }
            }).finally(function(){
                $timeout(function(){vm.ajax();}, 100000)
            });
        },

        more : function(index){
            vm.state[index] = state(vm.state[index]['v'] == 0 ? 1 : 0);
            vm.tsChange++;
        },

        option : function(which){
            vm[which + 'Demo'] = vm[which + 'Demo'] ? 0 : 1;
        }
    });
    vm.init();
    //vm.ajax();
}]);