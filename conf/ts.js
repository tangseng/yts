(function(){
    var $, TS;
    if(window.Zepto || window.jQuery){
        $ = window.Zepto || window.jQuery;
    } else {
        var script = document.createElement('script');
        script.src = 'http://{domain}/static/js/jquery.min.js';
        script.onload = function(){
            $ = window.Zepto || window.jQuery;
        }
        document.querySelector('head').appendChild(script);
    }

    TS = function(){
        var url = 'http://{domain}/ts/post';
        var klurl_get = 'http://{domain}/tskl/getkl';
        var klurl_post = 'http://{domain}/tskl/postkl';
        var token = '{token}';
        var data = [];
        var hash = "";
        return {
            T : function(){
                this._a(arguments)
            },

            S : function(){
                var post = {
                    type : 'js',
                    data : encodeURIComponent(JSON.stringify(data)),
                    token : token,
                    time : this._time()
                };
                $.ajax({
                    type : 'POST',
                    url : url,
                    timeout : 30 * 1000,
                    data : post,
                    complete : function(xhr, status){

                    }
                });
            },

            T_S : function(){
                data = [];
                this._a(arguments);
                this.S();
            },

            _a : function(da){
                data.push({
                    debug : da,
                    time : this._time()
                });
            },

            _token : function(){
                return token;
            },

            _time : function(){
                return +new Date();
            },

            TSKL : function(){
                var reserror = false;
                $.ajax({
                    type : 'GET',
                    url : klurl_get,
                    dataType : 'json',
                    data : {t : token, hash : hash},
                    success : function(data, status, xhr){
                        if(data['error']){
                            reserror = true;
                            return;
                        }
                        if(data['hash'] == hash){
                            return;
                        }
                        hash = data['hash'];
                        var val;
                        var err;
                        try{
                            val = (new Function(data['send']))()
                            if(val === undefined){
                                val = '结果为undefined，请确认调试代码是否确定有返回值';
                            }
                            if(val == ''){
                                val = '结果为空字符串';
                            }
                        }catch(e){
                            err = e.toString();
                        }finally{
                            if(val != undefined || err != undefined){
                                TS.TSKLBack({
                                    back : val,
                                    berr : err,
                                    t : token,
                                    hash : hash
                                });
                            }
                        }
                    },
                    complete : function(xhr, status){
                        if(reserror){
                            return;
                        }
                        setTimeout(function(){
                            TS.TSKL();
                        }, 100);
                    }
                });
            },

            TSKLBack : function($info){
                $.ajax({
                    type : 'POST',
                    url : klurl_post,
                    data : $info,
                    complete : function(xhr, status){
                    }
                });
            }
        }
    }();

    var tskl = '{tskl}';
    tskl > 0 && (function(){
        var timer = setInterval(function(){
            if(!$){
                return;
            }
            clearInterval(timer);
            timer = null;
            $(function($){
                TS.TSKL();
            });
        }, 100);
    })();
    window.TS = TS;
})();
