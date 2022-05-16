const App = {
    data() {
        return {
            message: "Hello Element Plus",
            api: '/api/todo',
            table: null,
            tableData: [],
            selected: [],
            form: {},
            dialogFormVisible: false,
        };
    },
    computed: {
        dialogFormTitle() { return this.form.id > 0 ? 'Update' : "Create" },
    },
    methods: {
        handleList() {
            axios.get(this.api)
                .then(response => (this.tableData = response.data.data))
                .catch(console.log)
        },
        handleCreate() {
            this.form = {}
            this.dialogFormVisible = true
        },
        handleUpdate(index, row) {
            this.form = row
            this.dialogFormVisible = true
        },
        handleSave() {
            this.dialogFormVisible = false
            if (this.form.id > 0) {
                axios.put(this.api, this.form)
                    .then(response => this.handleList())
                    .catch(console.log)
            } else {
                axios.post(this.api, this.form)
                    .then(response => this.handleList())
                    .catch(console.log)
            }
        },
        handleDelete(index, row) {
            this.handleMultiDelete([row])
        },
        handleMultiDelete(rows) {
            rows = rows || this.selected
            axios.delete(this.api, { params: { id: rows.map(x => x.id).join(',') } })
                .then(response => (this.tableData = response.data.data))
                .catch(console.log)
        },
    },
    mounted: function () {
        this.handleList()
    }
};

const app = Vue.createApp(App);
app.config.compilerOptions.delimiters = ['{${', '}}']
app.use(ElementPlus);
app.mount("#todo");
