<?php

class query
{
    protected $_sql;
    protected $_type;

    public function __construct($sql, $type = null)
	{
		$this->_sql = $sql;
		$this->_type = $type;
    }

    public function compile()
    {
        echo "query->compile() === \n";
    }

    public function execute()
    {
        // 这里调用的不一定是上面的compile()方法，select调用就是select的compile()方法，update、delete也是一样是他们对应的compile()方法
        // 好好反思一下golang的调用逻辑，明天见
        $this->compile();
        echo "query->execute() === \n";
    }
}

class builder extends query
{
}

class where extends builder
{
    protected $_where = array();
}

class select extends where
{
    protected $_select = array();

    public function __construct(array $columns = null)
	{
		if ( ! empty($columns))
		{
			$this->_select = $columns;
		}

		parent::__construct('', 1);
	}

    public function from()
	{
        echo "select->from() === \n";
		return $this;
    }

    public function compile()
    {
        $sql = $this->_sql;
        //echo "select === {$sql} \n";
        echo "select->compile() === \n";
    }
}

$s = new select();
//$s->from()->compile();
$s->from()->execute();



