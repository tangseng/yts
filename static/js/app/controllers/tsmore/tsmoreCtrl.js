angular.module('zfh').controller('ctrl.tsmore', ['$scope', '$http', '$timeout', "myTime", function($scope, $http, $timeout, myTime){
    var vm = $scope.vm = {};

    var state = function(bool){
        return {
            v : bool,
            n : bool ? '简单' : '详细'
        }
    };
    $.extend(vm, {
        tss : [],
        tsChange : 0,

        state : [],

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
            }
        },

        more : function(index){
            vm.state[index] = state(vm.state[index]['v'] == 0 ? 1 : 0);
            vm.tsChange++;
        }
    });
    vm.init();
}]);