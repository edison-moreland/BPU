module register_inc #(
    parameter N = 2
) (
    input  logic         write, read, inc,
    input  logic [N-1:0] data_in,
    output logic [N-1:0] data_out,
);

    logic [N-1:0] data;

    always_ff @(posedge write, posedge inc) begin
        if (write)
            data <= data_in;
        else if (inc)
            data <= data + 1;
    end

    always_comb begin
        if (read)
            data_out <= data;
        else
            data_out <= 0;
    end

endmodule