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
（1）initPaper稿件上链（具有相同标题title的拒绝上链）    
（2）getPaperInfo（通过标题得到文章信息）   
（3）getPapers（传入想要查询的json（例title：xxx）对进行富查询）   
（4）isPaperIllegal（通过paperHash检测，如果检测到相似或者相同的paperHash，就认为是一稿多投，这样把已经存在链上的文章信息返回出来；否则直接返回文章是合法的）   
（5）setPaperState（编辑通过文章标题title设置稿件状态 three退稿，直接从链上删除）   
