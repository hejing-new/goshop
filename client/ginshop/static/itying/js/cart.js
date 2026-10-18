(function($){
    var app={
        init:function(){
            this.changeCartNum();
            this.deleteConfirm();
            this.initCheckBox();
            this.isCheckedAll();
            this.initChekOut();
            this.updateCounts();        // ★ 首次渲染后校正"共 N 件 / 已选 M 件"
        },
        // ★ 金额统一保留两位小数，避免 3897.0000000000005 这类浮点噪声
        formatPrice:function(v){
            return (parseFloat(v) || 0).toFixed(2);
        },
        // ★ 按页面上的输入框和勾选状态实时统计，保证和表格里显示的数量一致
        updateCounts:function(){
            var total=0, checked=0;
            $(".cart_list").each(function(){
                var n = parseInt($(this).find(".input_center input").val(), 10) || 0;
                total += n;
                if ($(this).find(":checkbox").prop("checked")) {
                    checked += n;
                }
            });
            $("#totalNum").text(total);
            $("#checkedNum").text(checked);
        },
        deleteConfirm:function(){
            $('.delete').click(function(){
                return confirm('您确定要删除吗?');
            })
        },
        initChekOut(){
            $(function(){
                $("#checkout").click(function(){
                    var allPrice=parseFloat($("#allPrice").html());
                    if(!allPrice || allPrice==0){
                        alert('购物车没有选中去结算的商品')
                    }else{
                        location.href="/buy/checkout";
                    }
                })
            })
        },
        initCheckBox(){
            //全选 / 全不选
            $("#checkAll").click(function() {
                var flag = this.checked ? 1 : 0;
                $(":checkbox").prop("checked", this.checked);
                $.get('/cart/changeAllCart?flag='+flag,function(response){
                    if(response.success){
                        $("#allPrice").html(app.formatPrice(response.allPrice)+"元")
                        app.updateCounts();
                    }
                })
            });

            //点击单个选择框按钮的时候触发
            $(".cart_list :checkbox").click(function() {
                app.isCheckedAll();
                var goods_id=$(this).attr("goods_id")
                var goods_color=$(this).attr("goods_color")
                $.get('/cart/changeOneCart?goods_id='+goods_id+'&goods_color='+goods_color,function(response){
                    if(response.success){
                        $("#allPrice").html(app.formatPrice(response.allPrice)+"元")
                        app.updateCounts();
                    }
                })
            });
        },
        //判断全选是否选择
        isCheckedAll(){
            var allNum = $(".cart_list :checkbox").size();//checkbox总个数
            var checkedNum = 0;
            $(".cart_list :checkbox").each(function () {
                if($(this).prop("checked")==true){
                    checkedNum++;
                }
            });
            $("#checkAll").prop("checked", allNum==checkedNum);
        },
        changeCartNum(){
            $('.decCart').click(function(){
                var goods_id=$(this).attr("goods_id")
                var goods_color=$(this).attr("goods_color")
                var _that=this;
                $.get('/cart/decCart?goods_id='+goods_id+'&goods_color='+goods_color,function(response){
                    if(response.success){
                        $("#allPrice").html(app.formatPrice(response.allPrice)+"元")
                        $(_that).siblings(".input_center").find("input").val(response.num)
                        $(_that).parent().parent().siblings(".totalPrice").html(app.formatPrice(response.currentPrice)+"元")
                        app.updateCounts()
                        app.isCheckedAll()
                    }
                })
            });

            $('.incCart').click(function(){
                var goods_id=$(this).attr("goods_id")
                var goods_color=$(this).attr("goods_color")
                var _that=this;
                $.get('/cart/incCart?goods_id='+goods_id+'&goods_color='+goods_color,function(response){
                    if(response.success){
                        $("#allPrice").html(app.formatPrice(response.allPrice)+"元")
                        $(_that).siblings(".input_center").find("input").val(response.num)
                        $(_that).parent().parent().siblings(".totalPrice").html(app.formatPrice(response.currentPrice)+"元")
                        app.updateCounts()
                        app.isCheckedAll()
                    }
                })
            });
        }
    }

    $(function(){
        app.init();
    })
})($)
