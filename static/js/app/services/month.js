angular.module('zfh').factory('myMonth', [function(){
    function mm(m, y){
        if(m == 2){
            if(y % 4 == 0) {
                return 29;
            }
            return 28;
        }
        if($.inArray(m, [1,3,5,7,8,10,12]) != -1){
            return 31;
        }
        return 30;
    }

    function createDay(y, m, d, isToday) {
        var step = GlobalConfig.isSafari ? '/' : '-';
        var day = {
            'num' : d,
            'date' : y + step + pad(m) + step + pad(d),
            'y' : y,
            'm' : m,
            'd' : d,
            'today' : isToday,
            'active' : isToday == 1 ? true : false
        };
        day.dateStart = time(day.date + ' ' + GlobalConfig.start);
        day.dateEnd = time(day.date + ' ' + GlobalConfig.end);
        day.shangban = shangban(day.date);
        return day;
    }

    function shangban(date){
        var w = new Date(date).getDay();
        return  w != 0 && w != 6 ? 1 : 0;
    }

    function time (date){
        return +new Date(date) / 1000;
    }

    function pad(string, len){
        len = len || 2;
        string += '';
        return len - string.length ? (new Array(len - string.length + 1)).join("0") + string : string;
    }

    return function(Y, M, D){
        var days = [];
        if(D > 15){
            var cmm = mm(M, Y);
            for(var ii = 16; ii <= cmm; ii++){
                days.push(createDay(Y, M, ii, ii == D ? 1 : (ii > D ? 2 : 0)));
            }
            var nM = M == 12 ? 1 : M + 1;
            var nY = M == 12 ? Y + 1 : Y;
            for(ii = 1; ii <= 15; ii++){
                days.push(createDay(nY, nM, ii, 2));
            }
        } else {
            var pM = M == 1 ? 12 : M - 1;
            var pY = M == 1 ? Y -1 : Y;
            var cmm = mm(pM, Y);
            for(var ii = 16; ii <= cmm; ii++){
                days.push(createDay(pY, pM, ii, 0));
            }
            for(ii = 1; ii <= 15; ii++){
                days.push(createDay(Y, M, ii, ii == D ? 1 : (ii > D ? 2 : 0)));
            }
        }
        return days;
    }
}]);