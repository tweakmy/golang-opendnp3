package APL

import (
	//"fmt"
	"time"
)
import h "opendnp3/helper"

//Create a Logger object
type Logger struct {
	Name string   	//Specifify the topic of the Logger
	h.Observable 	//Inherit Observable from helper group
}

type LogSubcriber struct{
	Name string  	//Name of the subriber which can be Kafka, fmt, logfiles 
	h.Observer		
}

//This the custom Publish that calls the Log function
func (l *Logger) Log(v LogEntry){
	//l.Value = v
	l.Publish(v)
}

//Add a log subcriber
func (l *Logger) AddLogSubscriber(logSubscriber LogSubcriber) {
   l.AddObserver(logSubscriber)
}

//Delete a subcriber
func (l *Logger) RemoveLogSubscriber(nameOfSubscriber string) {
   
}

/******************************************************************/
type FilterLevel uint

const(
	LEV_EVENT FilterLevel =	0x01
	LEV_ERROR =		0x02
	LEV_WARNING =	0x04
	LEV_INFO  =		0x08
	LEV_INTERPRET =	0x10
	LEV_COMM =		0x20
	LEV_DEBUG =		0x40
)

type LogEntry struct{
	Time time.Time
	Message string
	DeviceName string
	Location string
	MFilterLevel FilterLevel
	ErrorCode int
}

/********************************************************************/
type LogStdio struct{
	LogSubcriber //Inherit LogSubriber
} 
//type Field struct {
//	Value int64
//	Observable
//}

//func (f *Field) Set(v int64){
//	f.Value = v
//	f.Publish(v)
//}
//
//
//func Listen(value interface{}){
//	fmt.Printf("new value 1: %v\n", value)
//}

//func Listen2(value interface{}){
//	fmt.Printf("new value 2: %v\n", value)
//}

//func main() {
//	v := &Field{}
//	v.AddObserver(ObserverFunc(Listen))
//	v.AddObserver(ObserverFunc(Listen2))
//	v.Set(105)
//	
//	fmt.Println("Hello, playground")
//}


//class EventLog : public ILogBase, private Uncopyable
//{
//public:
//
//	/** Immediate printing to minimize effect of debugging output on execution timing. */
//	//EventLog();
//	virtual ~EventLog();
//
//	Logger* GetLogger( FilterLevel aFilter, const std::string& aLoggerID );
//	Logger* GetExistingLogger( const std::string& aLoggerID );
//	void GetAllLoggers( std::vector<Logger*>& apLoggers);
//
//	/**
//	* Binds a listener to ALL log messages
//	*/
//	void AddLogSubscriber(ILogBase* apSubscriber);
//
//	/**
//	* Binds a listener to only certain error messages
//	*/
//	void AddLogSubscriber(ILogBase* apSubscriber, int aErrorCode);
//
//	/**
//	* Cancels a previous binding
//	*/
//	void RemoveLogSubscriber(ILogBase* apBase);
//
//	//implement the log function from ILogBase
//	void Log( const LogEntry& arEntry );
//	void SetVar(const std::string& aSource, const std::string& aVarName, int aValue);
//
//private:
//
//	bool SetContains(const std::set<int>& arSet, int aValue);
//
//	SigLock mLock;
//
//	//holds pointers to the loggers that have been distributed
//	typedef std::map<std::string, Logger*> LoggerMap;
//	LoggerMap mLogMap;
//	typedef std::map<ILogBase*, std::set<int> > SubscriberMap;
//	SubscriberMap mSubscribers;
//
//};

//	LogEntry(): mTime(TimeStamp::GetUTCTimeStamp()) {};
//
//	LogEntry( FilterLevel aLevel, const std::string& aDeviceName, const std::string& aLocation, const std::string& aMessage, int aErrorCode);
//
//	const std::string&	GetDeviceName() const {
//		return mDeviceName;
//	}
//	const std::string&	GetLocation() const {
//		return mLocation;
//	}
//	const std::string&	GetMessage() const {
//		return mMessage;
//	}
//	FilterLevel			GetFilterLevel() const {
//		return mFilterLevel;
//	}
//	std::string			GetTimeString() const {
//                return TimeStamp::UTCTimeStampToString(mTime);
//	}