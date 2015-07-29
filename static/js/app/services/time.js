angular.module('zfh').factory('myTime', [function(){
    return {
        time : function(date){
            return +new Date(date) / 1000;
        },

        timeToStr : function(time, hi){
            var step = GlobalConfig.isSafari ? '/' : '-';
            var dd = new Date(time);
            var ss = dd.getFullYear() + step + this.pad((dd.getMonth() + 1)) + step + this.pad(dd.getDate());
            if(!hi){
                return ss;
            }
            return ss + ' ' + this.pad(dd.getHours()) + ':' + this.pad(dd.getMinutes());
        },

        pad : function(string, len){
            len = len || 2;
            string += '';
            return len - string.length ? (new Array(len - string.length + 1)).join("0") + string : string;
        },

        cduration : function(start, end){
            var duration = end - start;
            if(end >= 13.5 * 3600 && start <= 12 * 3600){
                duration -= 1.5 * 3600;
            }
            return duration;
        },

        duration : function(duration){
            if(arguments.length == 2){
                duration = this.cduration(arguments[0], arguments[1]);
            }
            var hour = Math.floor(duration / 3600);
            var min = (duration / 60) % 60;
            return min ? hour + "时" + min + "分" : hour + "小时";
        },

        hm : function(time){
            return this.pad(Math.floor(time / 3600) + "", 2) + ':' + this.pad((time / 60) % 60 + "", 2);
        },

        h_m : function(time){
            return {
                h : this.pad(Math.floor(time / 3600) + "", 2),
                m : this.pad((time / 60) % 60 + "", 2)
            };
        },

        offsetTime : function(h, i){
            return parseInt(h) * 3600 + parseInt(i) * 60;
        }
    }
}]);