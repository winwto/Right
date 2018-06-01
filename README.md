# Right
基于区块链的一稿多投系统    
一、功能分析：   
1.智能合约功能    
paper的数据结构（      
  PaperHash string `json:"paperHash"` //文章hash        
  Title     string `json:"title"`     //题目    
	Author    string `json:"author"`    //作者    
	IsReview  string `json:"isReview"` //文章在审状态  1.文章存在没有在审one  2：文章在审two  3：文章退稿three 4：文章成功发表four   
	Reviewers string `json:"revewers"` //审稿机构   
	Tel       string `json:"tel"`      //联系方式   
  ）   
（1）initPaper稿件上链（具有相同标题title的拒绝上链）</br>
     {parm:("title", "author", "reviewers", "tel"  "paperHash")}</br>
     {return: (1)title在数据库中存在the title has been exist! (2)成功上链the paper onchain success!}</br>
（2）getPaperInfo（通过标题得到文章信息）   
     {parm:("tltle")}</br>
     {return: (1)title不存在No Paper! (2)查到信息paper的全部信息}</br>
（3）getPapers（传入想要查询的json（例title：xxx）对进行富查询）
     {parm:("查询K"，"查询V")}</br>
     {return:(1)title全部信息[paper1,paper2...]，当没有查到时[]}</br>
（4）isPaperIllegal（通过paperHash检测，如果检测到相似或者相同的paperHash，就认为是一稿多投，这样把已经存在链上的文章信息返回出来；否则直接返回文章是合法的）   
     {parm:("paperHash")}</br>
     {return:(1)是一稿多投Illegal+paper信息 (2)不是一稿多投the paper is legal}</br>
（5）setPaperState（编辑通过文章标题title设置稿件状态 three退稿，直接从链上删除）
     {parm:("title","isReview")}</br>
     {return:(1)title不存在the paper is not exist! (2)isReview=three delete paper success! (3)setPaperState success!!! (4)argument(2) must      be one of two、three、four!}
     
     
2.作者端功能</br>
 (1)right作者确权（将文章内容simHash之后上链）</br>
     {Parm：("title", "author", "reviewers(确权的时候不需要)", "tel", "paperHash(content->paperHash)")}</br>
     {Return:(1)、确权成功（paperHash、time）(2)、确权失败（该文章已经被确权）}</br>
 (2)submit作者投稿（将文章投给某个审稿机构）</br>
     
