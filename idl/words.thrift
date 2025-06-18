// Define the package name for the generated Go code as 'words'
namespace go words

// Include the 'base.thrift' file, which contains common response structures
include "base.thrift"

// New: Structure for pronunciation information, used to store pronunciation - related data of words
struct PronounceInfo {
    // 1 is the field sequence number, 'phonetic' stores the phonetic notation of the word
    1: string phonetic
    // 2 is the field sequence number, 'url' stores the link to the pronunciation audio file
    2: string url
}

// New: Structure for sentence information, used to store example sentences and their audio links
struct SentenceInfo {
    // 1 is the field sequence number, 'text' stores the text content of the example sentence
    1: string text
    // 2 is the field sequence number, 'audio_url' stores the link to the example sentence's audio file
    2: string audio_url
}

// Structure for word information, used to store detailed data of a word
struct Word {
    // 1 is the field sequence number, 'word_id' stores the unique identifier of the word
    1: i64 word_id
    // 2 is the field sequence number, 'word_name' stores the name of the word
    2: string word_name
    // 3 is the field sequence number, 'description' stores the description of the word
    3: string description
    // 4 is the field sequence number, 'explains' stores the definitions of the word
    4: string explains
    // 5 is the field sequence number, 'pronounce_us' stores the American pronunciation information. Changed from string to PronounceInfo type.
    5: PronounceInfo pronounce_us    
    // 6 is the field sequence number, 'pronounce_uk' stores the British pronunciation information. Changed from string to PronounceInfo type.
    6: PronounceInfo pronounce_uk    
    // 7 is the field sequence number, 'tag_id' stores the tag ID of the word. New field.
    7: i32 tag_id                    
    // 8 is the field sequence number, 'level' stores the current review level. New field.
    8: i32 level                     
    // 9 is the field sequence number, 'max_level' stores the maximum review level. New field.
    9: i32 max_level                 
    // 10 is the field sequence number, 'sentences' stores the list of example sentences. New field.
    10: list<SentenceInfo> sentences 
}

// Request structure for adding a new word
struct AddWordReq {
    // 1 is the field sequence number, 'user_id' stores the user ID. Not directly passed by the frontend, can be ignored.
    1: i64 user_id     // Not directly passed by the frontend, can be ignored
    // 2 is the field sequence number, 'word_name' stores the name of the word to be added, corresponding to the 'word' field in the HTTP request body
    2: string word_name (api.body="word")
    // 3 is the field sequence number, 'tag_id' stores the tag ID of the word, corresponding to the 'tag_id' field in the HTTP request body
    3: i32 tag_id (api.body="tag_id")
}

// Request structure for updating the word tag
struct UpdateWordTag {
    // 1 is the field sequence number, 'word_id' stores the ID of the word whose tag is to be updated, corresponding to the 'word_id' field in the HTTP request body
    1: i64 word_id (api.body="word_id")
    // 2 is the field sequence number, 'user_id' stores the user ID. Not directly passed by the frontend, can be ignored.
    2: i64 user_id    // Not directly passed by the frontend, can be ignored
    // 3 is the field sequence number, 'tag_id' stores the tag ID to be updated, corresponding to the 'tag_id' field in the HTTP request body
    3: i32 tag_id (api.body="tag_id")
}

// Request structure for getting the word list
struct WordListReq {
    // 1 is the field sequence number, 'user_id' stores the user ID
    1: i64 user_id
    // 2 is the field sequence number, 'offset' stores the pagination offset, corresponding to the 'offset' field in the HTTP query parameters
    2: i64 offset (api.query="offset")
    // 3 is the field sequence number, 'num' stores the number of words to be retrieved, corresponding to the 'num' field in the HTTP query parameters
    3: i64 num (api.query="num")
}

// Response structure for getting the word list
struct WordListResp {
    // 1 is the field sequence number, 'words_list' stores the retrieved list of words
    1: list <Word> words_list
    // 255 is the field sequence number, 'BaseResp' stores common response information from the 'base.thrift' file
    255: base.BaseResp BaseResp
}

// Request structure for the chat function
struct ChatReq {
    // 1 is the field sequence number, 'prompt' stores the chat prompt information, corresponding to the 'p' field in the HTTP query parameters
    1: string prompt (api.query="p")
    // 2 is the field sequence number, 'queryMsg' stores the chat query message, corresponding to the 'q' field in the HTTP query parameters
    2: string queryMsg (api.query="q")
    // 3 is the field sequence number, 'conversation_id' stores the conversation ID
    3: string conversation_id
}

// Response structure for the chat function
struct ChatResp {
    // 1 is the field sequence number, 'msg' stores the chat response message
    1: string msg
    // 2 is the field sequence number, 'extra' stores additional information in key - value pairs
    2: map<string,string> extra
}

// Request structure for getting the review word list
struct ReviewListReq {
    // 1 is the field sequence number, 'user_id' stores the user ID
    1: i64 user_id
}

// Structure for option items, used to store option information for review questions
struct OptionItem {
    // 1 is the field sequence number, 'description' stores the description of the option
    1: string description
    // 2 is the field sequence number, 'answer_list_id' stores the answer list ID corresponding to the option
    2: i64 answer_list_id 
}

// Structure for review questions, used to store detailed information about review questions
struct ReviewQuestion {
    // 1 is the field sequence number, 'question' stores the content of the review question
    1: string question 
    // 2 is the field sequence number, 'word_id' stores the ID of the word associated with the question
    2: i64 word_id
    // 3 is the field sequence number, 'question_type' stores the type of the question
    3: i64 question_type
    // 4 is the field sequence number, 'options' stores the list of options for the question
    4: list <OptionItem> options
    // 5 is the field sequence number, 'show_info' stores the information to be displayed for fill - in - the - blank questions. New optional field.
    5: optional list<string> show_info  // New field, used for displaying fill - in - the - blank questions
}

// Response structure for getting the review word list
struct ReviewListResp {
    // 1 is the field sequence number, 'total_num' stores the total number of review questions
    1: string total_num
    // 2 is the field sequence number, 'questions' stores the list of review questions
    2: list <ReviewQuestion> questions
}

// Request structure for submitting an answer
struct SubmitAnswerReq {
    // 1 is the field sequence number, 'user_id' stores the user ID
    1: i64 user_id 
    // 2 is the field sequence number, 'word_id' stores the ID of the word associated with the question, corresponding to the 'word_id' field in the HTTP request body
    2: i64 word_id (api.body="word_id")
    // 3 is the field sequence number, 'answer_id' stores the ID of the submitted answer, corresponding to the 'answer_id' field in the HTTP request body
    3: i64 answer_id (api.body="answer_id")
    // 4 is the field sequence number, 'question_type' stores the type of the question, corresponding to the 'question_type' field in the HTTP request body
    4: i64 question_type (api.body="question_type")
    // 5 is the field sequence number, 'filled_name' stores the answer for fill - in - the - blank questions, corresponding to the 'filled_name' field in the HTTP request body. New optional field.
    5: optional string filled_name (api.body="filled_name")  // New: Answer for fill - in - the - blank questions
}

// Response structure for submitting an answer
struct SubmitAnswerResp {
    // 1 is the field sequence number, 'is_correct' stores the flag indicating whether the answer is correct
    1: bool is_correct
    // 2 is the field sequence number, 'correct_answer_id' stores the ID of the correct answer
    2: i64 correct_answer_id
    // 3 is the field sequence number, 'message' stores the response message
    3: string message
    // 255 is the field sequence number, 'BaseResp' stores common response information from the 'base.thrift' file
    255: base.BaseResp BaseResp
}

// Request structure for getting the review progress
struct ReviewProgressReq {
    // 1 is the field sequence number, 'user_id' stores the user ID. Required field.
    1: required i64 user_id
}

// Response structure for getting the review progress
struct ReviewProgressResp {
    // 1 is the field sequence number, 'pending_review_count' stores the number of pending reviews
    1: i32 pending_review_count
    // 2 is the field sequence number, 'completed_review_count' stores the number of completed reviews
    2: i32 completed_review_count
    // 3 is the field sequence number, 'last_update_time' stores the last update time
    3: string last_update_time
    // 4 is the field sequence number, 'total_words_count' stores the total number of words
    4: i32 total_words_count        // Total number of words
    // 5 is the field sequence number, 'all_completed_count' stores the total number of words that have completed the review, retrieved from the 'review_progress' table
    5: i32 all_completed_count      // Total number of words that have completed the review (retrieved from the 'review_progress' table)
    // 255 is the field sequence number, 'BaseResp' stores common response information from the 'base.thrift' file
    255: base.BaseResp BaseResp
}

// Structure for tag information, used to store detailed information about word tags
struct TagInfo {
    // 1 is the field sequence number, 'tag_id' stores the unique identifier of the tag
    1: i64 tag_id
    // 2 is the field sequence number, 'tag_name' stores the name of the tag
    2: string tag_name
    // 3 is the field sequence number, 'question_types' stores the question types associated with the tag
    3: i32 question_types
    // 4 is the field sequence number, 'max_score' stores the maximum score corresponding to the tag
    4: i32 max_score
}

// Request structure for getting the supported tag types. No parameters are needed for now.
struct GetTagsReq {
    // No parameters are needed for now. Get all tag types.
}

// Response structure for getting the supported tag types
struct GetTagsResp {
    // 1 is the field sequence number, 'tags' stores the list of retrieved tag information
    1: list<TagInfo> tags
    // 255 is the field sequence number, 'BaseResp' stores common response information from the 'base.thrift' file
    255: base.BaseResp BaseResp
}

// Structure for word details, used to store detailed information about a word
struct WordDetail {
    // 1 is the field sequence number, 'new_word_name' stores the name of the word
    1: string new_word_name
    // 2 is the field sequence number, 'description' stores the description of the word
    2: string description
    // 3 is the field sequence number, 'explains' stores the definitions of the word
    3: string explains
    // 4 is the field sequence number, 'pronounce_us' stores the American pronunciation information
    4: PronounceInfo pronounce_us
    // 5 is the field sequence number, 'pronounce_uk' stores the British pronunciation information
    5: PronounceInfo pronounce_uk
    // 6 is the field sequence number, 'sentences' stores the list of example sentences
    6: list<SentenceInfo> sentences
}

// Response structure for word - related operations
struct WordResp {
    // 1 is the field sequence number, 'BaseResp' stores common response information from the 'base.thrift' file
    1: base.BaseResp BaseResp
    // 2 is the field sequence number, 'word' stores the word information. Optional field.
    2: optional Word word
}

// Define the word service, which contains multiple interfaces related to words
service WordService {
    // Chat function interface. Uses server - side streaming. Corresponds to the HTTP GET request '/api/chat'
    ChatResp Chat(1: ChatReq req) (streaming.mode="server", api.get="/api/chat")

    // Get the word list interface. Corresponds to the HTTP GET request '/api/word - list'
    WordListResp GetWordList(1: WordListReq req) (api.get="/api/word-list")
    // Get the word details by word name interface. Corresponds to the HTTP GET request '/api/word - detail'
    // Get the word details by word name
    WordDetail GetWordDetail(1: string word_name) (api.get="/api/word-detail")

    // Get the word information by word ID interface. Corresponds to the HTTP GET request '/api/word - query'
    WordResp GetWordByID(1: i64 word_id) (api.get="/api/word-query")
    // Status of adding a word:
    // 0    : Success. Returns the newly inserted word. 'error_msg' can be ignored.
    // 1    : Success, but the word already exists. No new insertion will be performed. Returns the existing word. 'error_msg' can be ignored.
    // -1   : Failure. The specific reason is returned via 'error_msg'. 'word' is empty.
    // Add a new word interface. Corresponds to the HTTP POST request '/api/word - add'
    WordResp AddNewWord (1: AddWordReq req) (api.post="/api/word-add")
    // 0    : Success. Returns the updated word. 'error_msg' can be ignored.
    // -1   : Failure. The specific reason is returned via 'error_msg'. 'word' is empty.
    // Update the word tag interface. Corresponds to the HTTP POST request '/api/word - tag'
    WordResp UpdateWordTagID(1: UpdateWordTag req) (api.post="/api/word-tag")

    // Get the review word list interface. Corresponds to the HTTP GET request '/api/review - list'
    ReviewListResp GetReviewWordList (1: ReviewListReq req) (api.get="/api/review-list")
    // Get today's review progress
    // Get today's review progress interface. Corresponds to the HTTP GET request '/api/review - progress'
    ReviewProgressResp GetTodayReviewProgress(1: ReviewProgressReq req) (api.get="/api/review-progress")

    // Submit an answer interface. Corresponds to the HTTP POST request '/api/answer'
    SubmitAnswerResp SubmitAnswer(1: SubmitAnswerReq req) (api.post="/api/answer")

    // New: Get the supported tag types
    // Get the supported tag types interface. Corresponds to the HTTP GET request '/api/tags'
    GetTagsResp GetSupportedTags(1: GetTagsReq req) (api.get="/api/tags")
}
