var GlobalConfig = {
    start : '09:00',
    end : '17:30',
    isSafari : navigator.userAgent.indexOf("Chrome") > -1 ? 0 : (navigator.userAgent.indexOf("Safari") > -1 ? 1 : 0)
};


angular.module('zfh', []).run(['$http', function($http){
    $http.defaults.headers.post["Content-Type"] = "application/x-www-form-urlencoded";
}]);
angular.element(document).ready(function() {
    angular.bootstrap(document, ['zfh']);
});


