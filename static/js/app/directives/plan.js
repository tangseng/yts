angular.module('zfh').directive('myPlan', ['myTime', function(myTime){
    return {
        restrict : "A",
        scope : {
            days : "=days",
            persons : "=persons"
        },
        link : function(scope, element, atts){
            var start = parseInt(atts['start']);
            var end = parseInt(atts['end']);
            var real = parseInt(atts['real']);
            var realS = real ? myTime.time(myTime.timeToStr(real * 1000) + ' ' + GlobalConfig.end) : 0;
            var color = atts['color'];
            var $this = element;
            var daySpans = [];
            $.each(scope.days, function(_, day){
                if(realS){
                    end = realS;
                }
                var cx = '';
                var bg = '';
                var border = '2px dotted ' + color + ';';
                if(day.dateStart >= start && day.dateEnd <= end){
                    bg += 'style="';
                    bg += 'border-top:' + border + 'border-bottom:' + border;
                    if(day.dateStart == start){
                        bg += 'border-left:' + border;
                    }
                    if(day.dateEnd == end){
                        bg += 'border-right:' + border;
                    }
                    if(day.today < 1){
                        bg += 'background:' + color + ';';
                    }
                    cx = 'll-se';
                    bg += '"';
                }

                daySpans.push('<span class="ll-day-mask ' + cx + '" ' + bg + ' data-day="' + day.num + '"></span>');
            });
            $this.html(daySpans.join(""));
            var ses = $this.find('.ll-se');
            var firstSE = ses.eq(0);
            var w = firstSE.outerWidth();
            var width = w * ses.length;
            var left = w * firstSE.prevAll('.ll-day-mask').length;
            var title = '【' + atts['status'] + '%】' + atts['title'];
            var bgSpan = atts['status'] == '100' && width > 0 ? '<span class="ll-plan-gou"></span>' : '';
            $('<a class="ll-plan-title" data-id="'+ atts['id'] +'" tabindex="0" data-trigger="focus" style="left:' + left + 'px;width:' + width + 'px;">' + bgSpan + '<span class="ll-plan-tt">' + title + '</span></a>').appendTo($this).popover({
                placement : 'bottom',
                html: true,
                title: (scope.persons[atts.uid] || {})['name'] + ' ：' + title,
                container: 'body',
                content: function(){
                    var $this = $(this);
                    var $tip = $this.data('bs.popover').$tip
                    $.get(
                        '/plan/getPlan',
                        {id : $this.data('id')},
                        function(data){
                            if(data && !data.error){
                                var html = '';
                                if(!data.length){
                                    html = '目前还没有日志';
                                } else {
                                    html += '<ul>';
                                    $.each(data, function(_, dd){
                                        html += '<li>';
                                        html += '<span>'+ myTime.timeToStr(dd.time * 1000, true) +' </span>';
                                        html += '<span>【'+ dd.status +'%】</span>';
                                        html += '<span>'+ dd.content +'</span>';
                                        html += '</li>';
                                    });
                                    html += '</ul>';
                                }
                                $tip.find('.popover-content').html(html);
                            }
                        }
                    );
                    return '<img src="/static/img/loading.gif" style="width:25px;"/>';
                }
            });
        }
    }
}]);