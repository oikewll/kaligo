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
        echo " === Object Name: " . (new ReflectionObject($this))->getName() . " === \n";
        echo "query->compile() === \n";
    }

    public function execute()
    {
        $this->compile();
        echo " === Object Name: " . (new ReflectionObject($this))->getName() . " === \n";
        echo "query->execute() === \n";
    }
}

class builder extends query
{
}

class where extends builder
{
    protected $_where = array();

    public function and_where($column, $op, $value)
    {
        $this->_where["AND"] = [$column, $op, $value];
        echo " === Object Name: " . (new ReflectionObject($this))->getName() . " === \n";
        return $this;
    }
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
$s->from()->and_where("uid", "=", "12")->execute();



