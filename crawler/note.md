# 大鱼号视频上传请求流程

1.  上传，URI地址(POST)：
        //这个地址的参数有两处是不一样的，它不是固定的也就是说需要我们自己去构造
        https://ums-origin-file.oss-cn-hangzhou.aliyuncs.com/video/wemedia/941878fe5eb1434a8f2afe13a3874c51/e274d55e-6cb4-11e8-b156-6c92bf56a17c?uploadId=480C8FD08DB74F8A96F121335903E4B9&callback=eyJjYWxsYmFja1VybCI6Imh0dHA6XC9cL2VuZ2luZS51bXMudWMuY25cL3ZpZGVvXC91cGxvYWRfbm90aWNlP2xvY2F0aW9uPW9zcy1jbi1oYW5nemhvdSZmaWxlX3R5cGU9dmlkZW8iLCJjYWxsYmFja0JvZHkiOiJvYmplY3Q9JHtvYmplY3R9JnNpemU9JHtzaXplfSZtaW1lVHlwZT0ke21pbWVUeXBlfSZidWNrZXQ9JHtidWNrZXR9JmV0YWc9JHtldGFnfSJ9
        https://ums-origin-file.oss-cn-hangzhou.aliyuncs.com/video/wemedia/941878fe5eb1434a8f2afe13a3874c51/108e046e-6cc0-11e8-8444-6c92bf56a17c?uploadId=A265828B8F7A43E09517EDCC0B6F69A0&callback=eyJjYWxsYmFja1VybCI6Imh0dHA6XC9cL2VuZ2luZS51bXMudWMuY25cL3ZpZGVvXC91cGxvYWRfbm90aWNlP2xvY2F0aW9uPW9zcy1jbi1oYW5nemhvdSZmaWxlX3R5cGU9dmlkZW8iLCJjYWxsYmFja0JvZHkiOiJvYmplY3Q9JHtvYmplY3R9JnNpemU9JHtzaXplfSZtaW1lVHlwZT0ke21pbWVUeXBlfSZidWNrZXQ9JHtidWNrZXR9JmV0YWc9JHtldGFnfSJ9
    POST数据：
        视频文件流
    响应：
        {"code":0,"message":"","data":{"id":"62a4f0ee497ccd8a"}}

2.  后台默认确认上传，URI地址(POST)：
        https://mp.dayu.com/confirm_info
    POST数据：
        videoId:62a4f0ee497ccd8a
        videoName:在水里吃饭睡觉看书的这种生活，敢不敢想象？.mp4
        videoDesc:
        appid:wemedia
        uid:941878fe5eb1434a8f2afe13a3874c51        //UID和utoken在user信息里
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"succeed":1}

3.  后台入库，URI地址(POST)：
        https://mp.dayu.com/dashboard/addMaterial
    POST数据：
        type:video
        title:在水里吃饭睡觉看书的这种生活，敢不敢想象？.mp4
        video[video_id]:62a4f0ee497ccd8a
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"data":
            { "_created_at":"2018-06-10T21:51:28.000+0800",                     "_id":"105d82d6386946228d51124d078ae334",
              "_updated_at":"2018-06-10T21:51:28.000+0800"
            }
        }
    
4. 生成视频播放地址,URi地址(POST):
        https://mp.dayu.com/dashboard/getVideoPlayUrl
    POST数据：
        id:62a4f0ee497ccd8a
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"playUrl":"https://v-ums.uc.cn/player/player.html?id=62a4f0ee497ccd8a&controls=1&pb=b&tm=1528643440&appid=wemedia&sign=e29cfb2b42b96989dc02cb3fdcbf1c5a"}

5. 生成封面，URI地址(GET):
        https://mp.dayu.com/media_keyframe?id=62a4f0ee497ccd8a&_=1528643270145
    响应(这里会有三次的请求，直到succeed为1也就是上传的视频加载完毕以后只有拿到状态为1的时候才有全部封面)：
        {"succeed":1,"api":"/1/images/result/new","data":{"id":"62a4f0ee497ccd8a","req_count":8,"resp_count":8,"total_count":8,"video_keyframes":{"frames":[{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/1","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":1},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/7","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":7},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/16","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":16},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/8","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":8},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/29","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":29},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/3","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":3},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/2","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":2},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/23","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":23}]}}}

6. 生成提交视频所有信息时的keyframe_index，URI地址(POST)
        https://mp.dayu.com/media_keyframe?id=62a4f0ee497ccd8a&_=1528643270145
    POST数据：
        id:62a4f0ee497ccd8a
        _:1528643270145
    响应：
        {"succeed":1,"api":"/1/images/result/new","data":{"id":"62a4f0ee497ccd8a","req_count":8,"resp_count":8,"total_count":8,"video_keyframes":{"frames":[{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/1","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":1},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/7","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":7},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/16","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":16},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/8","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":8},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/29","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":29},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/3","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":3},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/2","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":2},{"img_url":"http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/23","width":1920,"height":1080,"blackBorderPercent":0,"blackBorderLevel":"zero","blurMark":false,"blackBorderMark":false,"keyframe_index":23}]}}}

7. 选择封面，URI地址(POST)：
        https://ns.dayu.com/article/outsiteImgResave?appid=website&wmid=941878fe5eb1434a8f2afe13a3874c51&sign=edc37ce0f1938f9b91aff565f7e1c51b
    POST数据：
        imgSrc:http://img.ums.uc.cn/vkf/2-928448662-1334872304-203599805-0--1528643439732/29
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"code":0,"message":"","data":{"status":1,"url":"http://image.uc.cn/s/wemedia/s/upload/2018/3723da32099dbeabcea512436e6ee938x1777x1000x62.jpeg"}}

8. 保存封面，URI地址(POST):
        https://ns.dayu.com/article/imagecut?appid=website&wmid=941878fe5eb1434a8f2afe13a3874c51&sign=9338f1f93886fb374ec03654d51581f6
    POST数据：
        imgSrc:http://image.uc.cn/s/wemedia/s/upload/2018/3723da32099dbeabcea512436e6ee938x1777x1000x62.jpeg
        cutX:0
        cutY:0
        oriWidth:1770
        oriHeight:1000
        saveWidth:1024
        saveHeight:578
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"code":0,"message":"","data":{"status":1,"url":"http://image.uc.cn/s/wemedia/s/upload/2018/52252679c71859c492916a26c4e929c0x1024x578x29.jpeg"}}

9. 提交前的一部操作，URI地址(POST):
        https://mp.dayu.com/filterContent
    POST数据：
        author:
        title:在水里吃饭睡觉看书的这种生活，敢不敢想象？
        content:
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"result":0}

10. 提交视频的所有信息，URI地址(POST)：
        https://mp.dayu.com/dashboard/save-draft
    POST数据：
        keyframe_index:29       
        cover_from:3
        isCloseAdManual:false
        article_category:1
        draft_id:
        article_id:
        materialId:105d82d6386946228d51124d078ae334
        mid:
        videoId:62a4f0ee497ccd8a
        videoFileName:在水里吃饭睡觉看书的这种生活，敢不敢想象？.mp4
        title:在水里吃饭睡觉看书的这种生活，敢不敢想象？
        sub_title:
        coverImg:http://image.uc.cn/s/wemedia/s/upload/2018/52252679c71859c492916a26c4e929c0x1024x578x29.jpeg
        article_type:1
        category:娱乐
        aoyun:false
        weixin_promote:false
        curDaySubmit:false
        customize_tags:["水下生活"]
        is_show_ad:true
        is_video_ad_share:false
        is_original:0
        article_activity_title:
        article_activity_id:
        is_timed_release:false
        time_for_release:1528650469483
        vertical_cover_url:
        keyword:
        is_exclusive:false
        use_vertical_cover:false
        source:1
        originCoverImg:
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应：
        {"data":{"_created_at":"2018-06-10T23:29:25.000+0800","_id":"ac60de50c67a4638b8255b56b6a2a40e","_updated_at":"2018-06-10T23:29:25.000+0800","wm_update_at":"2018-06-10 23:29:25","previewUrl":"https://mobile.dayu.com/preview.html?type=1&wm_aid=ac60de50c67a4638b8255b56b6a2a40e&ts=1528644565979&acl=b549bca11bccc98f3ec1ff140fe7f7d2&uc_biz_str=S:custom%7CC:iflow_wm2&uc_param_str=frdnsnpfvecpntnwprdssskt"}}

11. 生成ID(好像这一步根本不需要), URI地址(GET):
        https://mp.dayu.com/dashboard/preview/getVideoPreviewSign
    响应:
        {"data":
            {"_wm_id":"941878fe5eb1434a8f2afe13a3874c51",
            "_timestamp":1528644566113,
            "_sign":"06372974bbfcf2e0a7d08cf9e21276ad"}
        }

12. 确认发布视频，URI地址(POST)：
        https://mp.dayu.com/dashboard/submit-article
    POST数据:
        dataList[0][_id]:ac60de50c67a4638b8255b56b6a2a40e
        dataList[0][isDraft]:1
        curDaySubmit:false
        utoken:1d24fa6122532c469e18cf838ef0838b
    响应(这里拿到响应)：

=======完===成===全===部===过===程=============================

