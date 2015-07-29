angular.module('zfh').controller('ctrl.person', ['$scope', '$http', '$timeout', function($scope, $http, $timeout){
    var vm = $scope.vm = {};

    $.extend(vm, {
        persons : [],
        ccperson : {},

        init : function(reset){
            vm.ccperson = {
                edit : 0,
                name : "",
                loginName : "",
                loginPass : ""
            };

            !reset && $.each(Persons, function(_, person){
                vm.persons.push(person);
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

        checkError : function(update){
            var error = false;
            $.each(vm.ccperson, function(k, v){
                if($.trim(v).length == 0){
                    error = vm.config["null"];
                    return false;
                }
            });
            !update && $.each(vm.persons, function(_, person){
                if(vm.ccperson.loginName == person.loginName){
                    error = vm.config["login"];
                    return false;
                }
            });
            return error;
        },

        doAdd : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/user/create', function(data){
                vm.persons.push($.extend({}, vm.ccperson, {loginName : data.loginName}));
            });
            return false;
        },

        update : function(loginName){
            $.each(vm.persons, function(_, person){
                if(person.loginName == loginName){
                    vm.ccperson = $.extend({}, person);
                }
            });
            vm.ccperson.edit = 1;
            return false;
        },

        doUpdate : function(){
            var error = vm.checkError(true);
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/user/update', function(data){
                $.each(vm.persons, function(index, person){
                    if(person.loginName == vm.ccperson.loginName){
                        vm.persons.splice(index, 1, $.extend({}, vm.ccperson));
                    }
                });
            });
            return false;
        },

        doAjax : function(url, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(vm.ccperson),
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

        doDelete : function(loginName){
            if(!confirm('确定要删除么？')){
                return false;
            }
            $http({
                method : "POST",
                url : '/user/delete',
                data : $.param({loginName : loginName}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    $.each(vm.persons, function(index, person){
                        if(person && person.loginName == loginName){
                            vm.persons.splice(index, 1);
                            return false;
                        }
                    });
                }
            });
            return false;
        },

        config : {
            "null" : "每项都需要填写",
            "login" : "登录名不能相同"
        }
    });

    vm.init();
}]);